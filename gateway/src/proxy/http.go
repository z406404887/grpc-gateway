package proxy

// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: addsvc.proto

/*
Package pb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/

import (
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"fmt"
	"encoding/json"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"strings"
)


type MyMux struct {
	conn *grpc.ClientConn
}

type URI struct {
	packageName string
	serviceName string
	version string
	method string
}

func (uri *URI) getServiceName() string {
	st := strings.Split(uri.serviceName, ".")
	serviceName := ""
	for _, v := range st {
		serviceName += strings.ToUpper(v[:1]) + v[1:]
	}
	return fmt.Sprintf("%v.%v", uri.packageName, serviceName)
}

func (uri *URI) getMethod() string {
	return strings.ToUpper(uri.method[:1]) + uri.method[1:]
}

func (p *MyMux) parseURL(url string) *URI {
	// /proto/service.add/v1/sum
	st := strings.Split(url, "/")
	if len(st) < 5 {
		return nil
	}
	return &URI{
		packageName:st[1],
		serviceName:st[2],
		version:st[3],
		method:st[4],
	}
}


func (p *MyMux) parseParams(req *http.Request) map[string]interface{} {
	req.ParseForm()
	//if strings.ToLower(req.Header.Get("Content-Type")) == "application/json" {
	// 处理传统意义上表单的参数，这里添加body内传输的json解析支持
	// 解析后的值默认追加到表单内部

	params := make(map[string]interface{})
	for key, v := range req.Form {
		params[key] = v[0]
	}
	if req.ContentLength <= 0 {
		return params
	}

	var data map[string]interface{}
	buf := make([]byte, req.ContentLength)
	n , err := req.Body.Read(buf)
	if err != nil || n <= 0 {
		fmt.Printf("req.Body read error: %v", err)
		return params
	}
	err = json.Unmarshal(buf, &data)
	if err != nil || data == nil {
		return params
	}
	for k, dv := range data {
		params[k] = dv
	}
	return params
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// url like:
	// http://localhost:8084/proto/service.add/v1/sum
	// package name is: proto
	// service name is: service.add
	// version is: v1
	// method is: sum
	fmt.Printf("%+v\n", *r)
	fmt.Printf("url: %+v\n", *r.URL)

	uri := p.parseURL(r.URL.Path)
	if uri == nil {
		w.Write([]byte("url path error, url path must be format by: /{packagename}/{servicename}/{version}/{method}"))
		return
	}
	fmt.Printf("uri: %+v\n", *uri)


	params := p.parseParams(r)//make(map[string] int)
	//params["a"] = 100
	//params["b"] = 200



	serviceName := uri.getServiceName()
	method := uri.getMethod()

	fullMethod := fmt.Sprintf("/%v/%v", serviceName, method)
	fmt.Printf("fullMethod=%s\v", fullMethod)

	var out interface{}
	err := p.conn.Invoke(context.Background(), fullMethod, params, &out)

	//buf := make([]byte, 1 << 20)
	//runtime.Stack(buf, true)
	//fmt.Printf("\n%s\n\n", buf)


	//debug.PrintStack()
	//out := new(SumReply)
	//err := grpc.Invoke(ctx, "/pb.Add/Sum", in, out, c.cc, opts...)
	//if err != nil {
	//	return nil, err
	//}
	//return out, nil
	fmt.Printf("return: %+v, error: %+v\n", out, err)
	b, _:=json.Marshal(out)
	w.Write(b)
	return
}

//var _ codes.Code
//var _ io.Reader
//var _ status.Status
//var _ = runtime.String
//var _ = utilities.NewDoubleArray
//
//func request_Add_Sum_0(ctx context.Context, marshaler runtime.Marshaler, client AddClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
//	var protoReq SumRequest
//	var metadata runtime.ServerMetadata
//
//	var (
//		val string
//		ok  bool
//		err error
//		_   = err
//	)
//
//	val, ok = pathParams["a"]
//	if !ok {
//		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "a")
//	}
//
//	protoReq.A, err = runtime.Int64(val)
//
//	if err != nil {
//		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "a", err)
//	}
//
//	val, ok = pathParams["b"]
//	if !ok {
//		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "b")
//	}
//
//	protoReq.B, err = runtime.Int64(val)
//
//	if err != nil {
//		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "b", err)
//	}
//
//	msg, err := client.Sum(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
//	return msg, metadata, err
//
//}


//type P struct{}


// RegisterAddHandlerFromEndpoint is same as RegisterAddHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterHandlerFromEndpoint(ctx context.Context, mux *MyMux, opts []grpc.DialOption) (err error) {
	address := "localhost:8082"//服务发现的地址
	opt1 := grpc.WithDefaultCallOptions(grpc.CallCustomCodec(MyCodec(encoding.GetCodec(proto.Name))))
	//opt2 := grpc.WithDefaultCallOptions(grpc.CallContentSubtype("proto"))
	opts = append(opts, opt1)
	//opts = append(opts, opt2)
	//grpc.NewContextWithServerTransportStream()
	mux.conn, err = grpc.Dial(address, opts...)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			if cerr := mux.conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", address, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := mux.conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", address, cerr)
			}
		}()
	}()

	return
}

// RegisterAddHandler registers the http handlers for service Add to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
//func RegisterHandler(ctx context.Context, mux *MyMux, conn *grpc.ClientConn) error {
//	return RegisterAddHandlerClient(ctx, mux, NewAddClient(conn))
//}

// RegisterAddHandler registers the http handlers for service Add to "mux".
// The handlers forward requests to the grpc endpoint over the given implementation of "AddClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "AddClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "AddClient" to call the correct interceptors.
//func RegisterAddHandlerClient(ctx context.Context, mux *runtime.ServeMux, client AddClient) error {
//	mux.Handle("GET", pattern_Add_Sum_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
//		ctx, cancel := context.WithCancel(req.Context())
//		defer cancel()
//		if cn, ok := w.(http.CloseNotifier); ok {
//			go func(done <-chan struct{}, closed <-chan bool) {
//				select {
//				case <-done:
//				case <-closed:
//					cancel()
//				}
//			}(ctx.Done(), cn.CloseNotify())
//		}
//		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
//		rctx, err := runtime.AnnotateContext(ctx, mux, req)
//		if err != nil {
//			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
//			return
//		}
//		resp, md, err := request_Add_Sum_0(rctx, inboundMarshaler, client, req, pathParams)
//		ctx = runtime.NewServerMetadataContext(ctx, md)
//		if err != nil {
//			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
//			return
//		}
//
//		forward_Add_Sum_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)
//
//	})
//
//	return nil
//}
//
//var (
//	pattern_Add_Sum_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 1, 0, 4, 1, 5, 1, 1, 0, 4, 1, 5, 2}, []string{"sum", "a", "b"}, ""))
//)
//
//var (
//	forward_Add_Sum_0 = runtime.ForwardResponseMessage
//)
