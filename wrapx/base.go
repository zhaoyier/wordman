package wrapx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"unicode"

	// "git.ezbuy.me/ezbuy/evtalk/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (b *_Base) register(router gin.IRouter, cList ...interface{}) bool {
	// groupPath := b.BasePath(router)
	mp := getInfo()
	for _, c := range cList {
		refTyp := reflect.TypeOf(c)
		refVal := reflect.ValueOf(c)
		t := reflect.Indirect(refVal).Type()
		objName := t.Name()

		// Install the methods
		for m := 0; m < refTyp.NumMethod(); m++ {
			method := refTyp.Method(m)
			num, _b := b.checkHandlerFunc(method.Type /*.Interface()*/, true)
			if _b {
				if v, ok := mp[objName+"."+method.Name]; ok {
					for _, v1 := range v {
						b.registerHandlerObj(router, v1.GenComment.Methods, v1.GenComment.RouterPath, method.Name, method.Func, refVal)
					}
				} else { // not find using default case
					routerPath, methods := b.getDefaultComments(objName, method.Name, num)
					fmt.Printf("=====>>001:%+v|%+v|%+v\n", objName, routerPath, methods)
					b.registerHandlerObj(router, methods, routerPath, method.Name, method.Func, refVal)
				}
			}
		}
	}
	return true
}

// checkHandlerFunc Judge whether to match rules
func (b *_Base) checkHandlerFunc(typ reflect.Type, isObj bool) (int, bool) { // 判断是否匹配规则,返回参数个数
	offset := 0
	if isObj {
		offset = 1
	}
	num := typ.NumIn() - offset
	if num == 1 || num == 2 { // Parameter checking 参数检查
		ctxType := typ.In(0 + offset)

		// go-gin default method
		if ctxType == reflect.TypeOf(&gin.Context{}) {
			return num, true
		}

		// Customized context . 自定义的context
		if ctxType == b.apiType {
			return num, true
		}

		// maybe interface
		if b.apiType.ConvertibleTo(ctxType) {
			return num, true
		}

	}
	return num, false
}

// registerHandlerObj Multiple registration methods.获取并过滤要绑定的参数
func (b *_Base) registerHandlerObj(router gin.IRouter, httpMethod []string, relativePath, methodName string, tvl, obj reflect.Value) error {
	call := b.handlerFuncObj(tvl, obj, methodName)

	for _, v := range httpMethod {
		// method := strings.ToUpper(v)
		// switch method{
		// case "ANY":
		// 	router.Any(relativePath,list...)
		// default:
		// 	router.Handle(method,relativePath,list...)
		// }
		// or
		switch strings.ToUpper(v) {
		case "POST":
			router.POST(relativePath, call)
		case "GET":
			router.GET(relativePath, call)
		case "DELETE":
			router.DELETE(relativePath, call)
		case "PATCH":
			router.PATCH(relativePath, call)
		case "PUT":
			router.PUT(relativePath, call)
		case "OPTIONS":
			router.OPTIONS(relativePath, call)
		case "HEAD":
			router.HEAD(relativePath, call)
		case "ANY":
			router.Any(relativePath, call)
		default:
			return fmt.Errorf("method:[%v] not support", httpMethod)
		}
	}

	return nil
}

func (b *_Base) getDefaultComments(objName, objFunc string, num int) (routerPath string, methods []string) {
	methods = []string{"ANY"}
	if num == 2 { // parm 2 , post default
		methods = []string{"post"}
	}

	if b.isBigCamel { // big camel style.大驼峰
		routerPath = b.prefix + "/" + b.service + "." + objName + "/" + objFunc
	} else {
		routerPath = b.prefix + "/" + b.service + "." + Ucfirst(objName) + "/" + Ucfirst(objFunc)
	}

	return
}

//首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// HandlerFunc Get and filter the parameters to be bound (object call type)
func (b *_Base) handlerFuncObj(tvl, obj reflect.Value, methodName string) gin.HandlerFunc { // 获取并过滤要绑定的参数(obj 对象类型)
	typ := tvl.Type()
	if typ.NumIn() == 2 { // Parameter checking 参数检查
		ctxType := typ.In(1)

		// go-gin default method
		apiFun := func(c *gin.Context) interface{} { return c }
		if ctxType == b.apiType { // Customized context . 自定义的context
			apiFun = b.apiFun
		} else if !(ctxType == reflect.TypeOf(&gin.Context{})) {
			panic("method " + runtime.FuncForPC(tvl.Pointer()).Name() + " not support!")
		}

		return func(c *gin.Context) {
			tvl.Call([]reflect.Value{obj, reflect.ValueOf(apiFun(c))})
		}
	}

	// Custom context type with request parameters .自定义的context类型,带request 请求参数
	call, err := b.getCallObj3(tvl, obj, methodName)
	if err != nil { // Direct reporting error.
		panic(err)
	}

	return call
}

// Custom context type with request parameters
func (b *_Base) getCallObj3(tvl, obj reflect.Value, methodName string) (func(*gin.Context), error) {
	typ := tvl.Type()
	if typ.NumIn() != 3 { // Parameter checking 参数检查
		return nil, errors.New("method " + runtime.FuncForPC(tvl.Pointer()).Name() + " not support!")
	}

	if typ.NumOut() != 0 {
		if typ.NumOut() == 2 { // Parameter checking 参数检查
			if returnType := typ.Out(1); returnType != typeOfError {
				return nil, fmt.Errorf("method : %v , returns[1] %v not error",
					runtime.FuncForPC(tvl.Pointer()).Name(), returnType.String())
			}
		} else {
			return nil, fmt.Errorf("method : %v , Only 2 return values (obj, error) are supported", runtime.FuncForPC(tvl.Pointer()).Name())
		}
	}

	ctxType, reqType := typ.In(1), typ.In(2)

	reqIsGinCtx := false
	if ctxType == reflect.TypeOf(&gin.Context{}) {
		reqIsGinCtx = true
	}

	// ctxType != reflect.TypeOf(gin.Context{}) &&
	// ctxType != reflect.Indirect(reflect.ValueOf(b.iAPIType)).Type()
	if !reqIsGinCtx && ctxType != b.apiType && !b.apiType.ConvertibleTo(ctxType) {
		return nil, errors.New("method " + runtime.FuncForPC(tvl.Pointer()).Name() + " first parm not support!")
	}

	reqIsValue := true
	if reqType.Kind() == reflect.Ptr {
		reqIsValue = false
	}
	apiFun := func(c *gin.Context) interface{} { return c }
	if !reqIsGinCtx {
		apiFun = b.apiFun
	}

	return func(c *gin.Context) {
		req := reflect.New(reqType)
		if !reqIsValue {
			req = reflect.New(reqType.Elem())
		}
		if err := b.unmarshal(c, req.Interface()); err != nil { // Return error message.返回错误信息
			b.handErrorString(c, req, err)
			return
		}

		if reqIsValue {
			req = req.Elem()
		}

		bainfo, is := b.beforCall(c, tvl, obj, req.Interface(), methodName)
		if !is {
			c.JSON(http.StatusBadRequest, bainfo.Resp)
			return
		}

		var returnValues []reflect.Value
		returnValues = tvl.Call([]reflect.Value{obj, reflect.ValueOf(apiFun(c)), req})

		if returnValues != nil {
			bainfo.Resp = returnValues[0].Interface()
			rerr := returnValues[1].Interface()
			if rerr != nil {
				bainfo.Error = rerr.(error)
			}

			is = b.afterCall(bainfo, obj)
			if is {
				c.JSON(http.StatusOK, bainfo.Resp)
			} else {
				c.JSON(http.StatusBadRequest, bainfo.Resp)
			}
		}
	}, nil
}

func (b *_Base) afterCall(info *GinBeforeAfterInfo, obj reflect.Value) bool {
	is := true
	if bfobj, ok := obj.Interface().(GinBeforeAfter); ok { // 本类型
		is = bfobj.GinAfter(info)
	}
	if is && b.beforeAfter != nil {
		is = b.beforeAfter.GinAfter(info)
	}
	return is
}

func (b *_Base) unmarshal(c *gin.Context, v interface{}) error {
	return c.ShouldBind(v)
	return nil
}

func (b *_Base) handErrorString(c *gin.Context, req reflect.Value, err error) {
	var fields []string
	if _, ok := err.(validator.ValidationErrors); ok {
		for _, err := range err.(validator.ValidationErrors) {
			tmp := fmt.Sprintf("%v:%v", b.FindTag(req.Interface(), err.Field(), "json"), err.Tag())
			if len(err.Param()) > 0 {
				tmp += fmt.Sprintf("[%v](but[%v])", err.Param(), err.Value())
			}
			fields = append(fields, tmp)
			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			// fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()
		}
	} else if _, ok := err.(*json.UnmarshalTypeError); ok {
		err := err.(*json.UnmarshalTypeError)
		tmp := fmt.Sprintf("%v:%v(but[%v])", err.Field, err.Type.String(), err.Value)
		fields = append(fields, tmp)

	} else {
		fields = append(fields, err.Error())
	}

	msg := b.GetErrorMsg(ParameterInvalid)
	msg.Error = fmt.Sprintf("req param : %v", strings.Join(fields, ";"))
	c.JSON(http.StatusBadRequest, msg)
	return
}

func (b *_Base) beforCall(c *gin.Context, tvl, obj reflect.Value, req interface{}, methodName string) (*GinBeforeAfterInfo, bool) {
	info := &GinBeforeAfterInfo{
		C:        c,
		FuncName: fmt.Sprintf("%v.%v", reflect.Indirect(obj).Type().Name(), methodName), // 函数名
		Req:      req,                                                                   // 调用前的请求参数
		Context:  context.Background(),                                                  // 占位参数，可用于存储其他参数，前后连接可用
	}

	is := true
	if bfobj, ok := obj.Interface().(GinBeforeAfter); ok { // 本类型
		is = bfobj.GinBefore(info)
	}
	if is && b.beforeAfter != nil {
		is = b.beforeAfter.GinBefore(info)
	}
	return info, is
}

// FindTag find struct of tag string.查找struct 的tag信息
func (b *_Base) FindTag(obj interface{}, field, tag string) string {
	dataStructType := reflect.Indirect(reflect.ValueOf(obj)).Type()
	for i := 0; i < dataStructType.NumField(); i++ {
		fd := dataStructType.Field(i)
		if fd.Name == field {
			bb := fd.Tag
			sqlTag := bb.Get(tag)

			if sqlTag == "-" || bb == "-" {
				return ""
			}

			sqlTags := strings.Split(sqlTag, ",")
			sqlFieldName := fd.Name // default
			if len(sqlTags[0]) > 0 {
				sqlFieldName = sqlTags[0]
			}
			return sqlFieldName
		}
	}

	return ""
}

//GetErrorMsg 获取错误消息 参数(int,string)
func (b *_Base) GetErrorMsg(errorCode ...interface{}) (msg MessageBody) {
	if len(errorCode) == 0 {
		fmt.Errorf("unknow")
		msg.State = false
		msg.Code = -1
		return
	}
	msg.State = false
	for _, e := range errorCode {
		switch v := e.(type) {
		case int:
			msg.Code = int(v)
			msg.Error = ErrCode(v).String()
		case ErrCode:
			_tryRegisteryCode(v)
			msg.Code = int(v)
			msg.Error = v.String()
		case string:
			msg.Error = string(v)
		case error:
			msg.Error = v.Error()
		default:
			msg.Error = fmt.Sprintf("Unknow type:(%v)", e)
		}
	}
	return
}

// Custom context type with request parameters
func (b *_Base) getCallFunc3(tvl reflect.Value) (func(*gin.Context), error) {
	typ := tvl.Type()
	if typ.NumIn() != 2 { // Parameter checking 参数检查
		return nil, errors.New("method " + runtime.FuncForPC(tvl.Pointer()).Name() + " not support!")
	}

	if typ.NumOut() != 0 {
		if typ.NumOut() == 2 { // Parameter checking 参数检查
			if returnType := typ.Out(1); returnType != typeOfError {
				return nil, fmt.Errorf("method : %v , returns[1] %v not error",
					runtime.FuncForPC(tvl.Pointer()).Name(), returnType.String())
			}
		} else {
			return nil, fmt.Errorf("method : %v , Only 2 return values (obj, error) are supported", runtime.FuncForPC(tvl.Pointer()).Name())
		}
	}

	ctxType, reqType := typ.In(0), typ.In(1)

	reqIsGinCtx := false
	if ctxType == reflect.TypeOf(&gin.Context{}) {
		reqIsGinCtx = true
	}

	// ctxType != reflect.TypeOf(gin.Context{}) &&
	// ctxType != reflect.Indirect(reflect.ValueOf(b.iAPIType)).Type()
	if !reqIsGinCtx && ctxType != b.apiType && !b.apiType.ConvertibleTo(ctxType) {
		return nil, errors.New("method " + runtime.FuncForPC(tvl.Pointer()).Name() + " first parm not support!")
	}

	reqIsValue := true
	if reqType.Kind() == reflect.Ptr {
		reqIsValue = false
	}
	apiFun := func(c *gin.Context) interface{} { return c }
	if !reqIsGinCtx {
		apiFun = b.apiFun
	}

	return func(c *gin.Context) {
		req := reflect.New(reqType)
		if !reqIsValue {
			req = reflect.New(reqType.Elem())
		}
		if err := b.unmarshal(c, req.Interface()); err != nil { // Return error message.返回错误信息
			b.handErrorString(c, req, err)
			return
		}

		if reqIsValue {
			req = req.Elem()
		}
		var returnValues []reflect.Value
		returnValues = tvl.Call([]reflect.Value{reflect.ValueOf(apiFun(c)), req})

		if returnValues != nil {
			obj := returnValues[0].Interface()
			rerr := returnValues[1].Interface()
			if rerr != nil {
				msg := b.GetErrorMsg(InValidOp)
				msg.Error = rerr.(error).Error()
				c.JSON(http.StatusBadRequest, msg)
			} else {
				c.JSON(http.StatusOK, obj)
			}
		}
	}, nil
}
