// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package blobbergrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// BlobberClient is the client API for Blobber service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlobberClient interface {
	GetAllocation(ctx context.Context, in *GetAllocationRequest, opts ...grpc.CallOption) (*GetAllocationResponse, error)
	GetFileMetaData(ctx context.Context, in *GetFileMetaDataRequest, opts ...grpc.CallOption) (*GetFileMetaDataResponse, error)
	GetFileStats(ctx context.Context, in *GetFileStatsRequest, opts ...grpc.CallOption) (*GetFileStatsResponse, error)
	ListEntities(ctx context.Context, in *ListEntitiesRequest, opts ...grpc.CallOption) (*ListEntitiesResponse, error)
	GetObjectPath(ctx context.Context, in *GetObjectPathRequest, opts ...grpc.CallOption) (*GetObjectPathResponse, error)
	GetReferencePath(ctx context.Context, in *GetReferencePathRequest, opts ...grpc.CallOption) (*GetReferencePathResponse, error)
	GetObjectTree(ctx context.Context, in *GetObjectTreeRequest, opts ...grpc.CallOption) (*GetObjectTreeResponse, error)
	Commit(ctx context.Context, in *CommitRequest, opts ...grpc.CallOption) (*CommitResponse, error)
	CalculateHash(ctx context.Context, in *CalculateHashRequest, opts ...grpc.CallOption) (*CalculateHashResponse, error)
	CommitMetaTxn(ctx context.Context, in *CommitMetaTxnRequest, opts ...grpc.CallOption) (*CommitMetaTxnResponse, error)
}

type blobberClient struct {
	cc grpc.ClientConnInterface
}

func NewBlobberClient(cc grpc.ClientConnInterface) BlobberClient {
	return &blobberClient{cc}
}

func (c *blobberClient) GetAllocation(ctx context.Context, in *GetAllocationRequest, opts ...grpc.CallOption) (*GetAllocationResponse, error) {
	out := new(GetAllocationResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/GetAllocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) GetFileMetaData(ctx context.Context, in *GetFileMetaDataRequest, opts ...grpc.CallOption) (*GetFileMetaDataResponse, error) {
	out := new(GetFileMetaDataResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/GetFileMetaData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) GetFileStats(ctx context.Context, in *GetFileStatsRequest, opts ...grpc.CallOption) (*GetFileStatsResponse, error) {
	out := new(GetFileStatsResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/GetFileStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) ListEntities(ctx context.Context, in *ListEntitiesRequest, opts ...grpc.CallOption) (*ListEntitiesResponse, error) {
	out := new(ListEntitiesResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/ListEntities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) GetObjectPath(ctx context.Context, in *GetObjectPathRequest, opts ...grpc.CallOption) (*GetObjectPathResponse, error) {
	out := new(GetObjectPathResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/GetObjectPath", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) GetReferencePath(ctx context.Context, in *GetReferencePathRequest, opts ...grpc.CallOption) (*GetReferencePathResponse, error) {
	out := new(GetReferencePathResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/GetReferencePath", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) GetObjectTree(ctx context.Context, in *GetObjectTreeRequest, opts ...grpc.CallOption) (*GetObjectTreeResponse, error) {
	out := new(GetObjectTreeResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/GetObjectTree", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) Commit(ctx context.Context, in *CommitRequest, opts ...grpc.CallOption) (*CommitResponse, error) {
	out := new(CommitResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/Commit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) CalculateHash(ctx context.Context, in *CalculateHashRequest, opts ...grpc.CallOption) (*CalculateHashResponse, error) {
	out := new(CalculateHashResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/CalculateHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobberClient) CommitMetaTxn(ctx context.Context, in *CommitMetaTxnRequest, opts ...grpc.CallOption) (*CommitMetaTxnResponse, error) {
	out := new(CommitMetaTxnResponse)
	err := c.cc.Invoke(ctx, "/blobber.service.v1.Blobber/CommitMetaTxn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlobberServer is the server API for Blobber service.
// All implementations must embed UnimplementedBlobberServer
// for forward compatibility
type BlobberServer interface {
	GetAllocation(context.Context, *GetAllocationRequest) (*GetAllocationResponse, error)
	GetFileMetaData(context.Context, *GetFileMetaDataRequest) (*GetFileMetaDataResponse, error)
	GetFileStats(context.Context, *GetFileStatsRequest) (*GetFileStatsResponse, error)
	ListEntities(context.Context, *ListEntitiesRequest) (*ListEntitiesResponse, error)
	GetObjectPath(context.Context, *GetObjectPathRequest) (*GetObjectPathResponse, error)
	GetReferencePath(context.Context, *GetReferencePathRequest) (*GetReferencePathResponse, error)
	GetObjectTree(context.Context, *GetObjectTreeRequest) (*GetObjectTreeResponse, error)
	Commit(context.Context, *CommitRequest) (*CommitResponse, error)
	CalculateHash(context.Context, *CalculateHashRequest) (*CalculateHashResponse, error)
	CommitMetaTxn(context.Context, *CommitMetaTxnRequest) (*CommitMetaTxnResponse, error)
	mustEmbedUnimplementedBlobberServer()
}

// UnimplementedBlobberServer must be embedded to have forward compatible implementations.
type UnimplementedBlobberServer struct {
}

func (UnimplementedBlobberServer) GetAllocation(context.Context, *GetAllocationRequest) (*GetAllocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllocation not implemented")
}
func (UnimplementedBlobberServer) GetFileMetaData(context.Context, *GetFileMetaDataRequest) (*GetFileMetaDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileMetaData not implemented")
}
func (UnimplementedBlobberServer) GetFileStats(context.Context, *GetFileStatsRequest) (*GetFileStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileStats not implemented")
}
func (UnimplementedBlobberServer) ListEntities(context.Context, *ListEntitiesRequest) (*ListEntitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEntities not implemented")
}
func (UnimplementedBlobberServer) GetObjectPath(context.Context, *GetObjectPathRequest) (*GetObjectPathResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObjectPath not implemented")
}
func (UnimplementedBlobberServer) GetReferencePath(context.Context, *GetReferencePathRequest) (*GetReferencePathResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReferencePath not implemented")
}
func (UnimplementedBlobberServer) GetObjectTree(context.Context, *GetObjectTreeRequest) (*GetObjectTreeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObjectTree not implemented")
}
func (UnimplementedBlobberServer) Commit(context.Context, *CommitRequest) (*CommitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedBlobberServer) CalculateHash(context.Context, *CalculateHashRequest) (*CalculateHashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateHash not implemented")
}
func (UnimplementedBlobberServer) CommitMetaTxn(context.Context, *CommitMetaTxnRequest) (*CommitMetaTxnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitMetaTxn not implemented")
}
func (UnimplementedBlobberServer) mustEmbedUnimplementedBlobberServer() {}

// UnsafeBlobberServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlobberServer will
// result in compilation errors.
type UnsafeBlobberServer interface {
	mustEmbedUnimplementedBlobberServer()
}

func RegisterBlobberServer(s *grpc.Server, srv BlobberServer) {
	s.RegisterService(&_Blobber_serviceDesc, srv)
}

func _Blobber_GetAllocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).GetAllocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/GetAllocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).GetAllocation(ctx, req.(*GetAllocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_GetFileMetaData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileMetaDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).GetFileMetaData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/GetFileMetaData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).GetFileMetaData(ctx, req.(*GetFileMetaDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_GetFileStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).GetFileStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/GetFileStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).GetFileStats(ctx, req.(*GetFileStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_ListEntities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEntitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).ListEntities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/ListEntities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).ListEntities(ctx, req.(*ListEntitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_GetObjectPath_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetObjectPathRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).GetObjectPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/GetObjectPath",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).GetObjectPath(ctx, req.(*GetObjectPathRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_GetReferencePath_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReferencePathRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).GetReferencePath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/GetReferencePath",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).GetReferencePath(ctx, req.(*GetReferencePathRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_GetObjectTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetObjectTreeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).GetObjectTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/GetObjectTree",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).GetObjectTree(ctx, req.(*GetObjectTreeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/Commit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).Commit(ctx, req.(*CommitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_CalculateHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalculateHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).CalculateHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/CalculateHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).CalculateHash(ctx, req.(*CalculateHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blobber_CommitMetaTxn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitMetaTxnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobberServer).CommitMetaTxn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blobber.service.v1.Blobber/CommitMetaTxn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobberServer).CommitMetaTxn(ctx, req.(*CommitMetaTxnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Blobber_serviceDesc = grpc.ServiceDesc{
	ServiceName: "blobber.service.v1.Blobber",
	HandlerType: (*BlobberServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllocation",
			Handler:    _Blobber_GetAllocation_Handler,
		},
		{
			MethodName: "GetFileMetaData",
			Handler:    _Blobber_GetFileMetaData_Handler,
		},
		{
			MethodName: "GetFileStats",
			Handler:    _Blobber_GetFileStats_Handler,
		},
		{
			MethodName: "ListEntities",
			Handler:    _Blobber_ListEntities_Handler,
		},
		{
			MethodName: "GetObjectPath",
			Handler:    _Blobber_GetObjectPath_Handler,
		},
		{
			MethodName: "GetReferencePath",
			Handler:    _Blobber_GetReferencePath_Handler,
		},
		{
			MethodName: "GetObjectTree",
			Handler:    _Blobber_GetObjectTree_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _Blobber_Commit_Handler,
		},
		{
			MethodName: "CalculateHash",
			Handler:    _Blobber_CalculateHash_Handler,
		},
		{
			MethodName: "CommitMetaTxn",
			Handler:    _Blobber_CommitMetaTxn_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blobber.proto",
}
