// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: tradesvc/v1/tradesvc.proto

package tradesvc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Side int32

const (
	Side_SIDE_UNSPECIFIED Side = 0
	Side_SIDE_BUY         Side = 1
	Side_SIDE_SELL        Side = 2
)

// Enum value maps for Side.
var (
	Side_name = map[int32]string{
		0: "SIDE_UNSPECIFIED",
		1: "SIDE_BUY",
		2: "SIDE_SELL",
	}
	Side_value = map[string]int32{
		"SIDE_UNSPECIFIED": 0,
		"SIDE_BUY":         1,
		"SIDE_SELL":        2,
	}
)

func (x Side) Enum() *Side {
	p := new(Side)
	*p = x
	return p
}

func (x Side) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Side) Descriptor() protoreflect.EnumDescriptor {
	return file_tradesvc_v1_tradesvc_proto_enumTypes[0].Descriptor()
}

func (Side) Type() protoreflect.EnumType {
	return &file_tradesvc_v1_tradesvc_proto_enumTypes[0]
}

func (x Side) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Side.Descriptor instead.
func (Side) EnumDescriptor() ([]byte, []int) {
	return file_tradesvc_v1_tradesvc_proto_rawDescGZIP(), []int{0}
}

type OrderType int32

const (
	OrderType_ORDER_TYPE_UNSPECIFIED   OrderType = 0
	OrderType_ORDER_TYPE_MARKET        OrderType = 1
	OrderType_ORDER_TYPE_LIMIT         OrderType = 2
	OrderType_ORDER_TYPE_STOP          OrderType = 3
	OrderType_ORDER_TYPE_STOP_LIMIT    OrderType = 4
	OrderType_ORDER_TYPE_TRAILING_STOP OrderType = 5
)

// Enum value maps for OrderType.
var (
	OrderType_name = map[int32]string{
		0: "ORDER_TYPE_UNSPECIFIED",
		1: "ORDER_TYPE_MARKET",
		2: "ORDER_TYPE_LIMIT",
		3: "ORDER_TYPE_STOP",
		4: "ORDER_TYPE_STOP_LIMIT",
		5: "ORDER_TYPE_TRAILING_STOP",
	}
	OrderType_value = map[string]int32{
		"ORDER_TYPE_UNSPECIFIED":   0,
		"ORDER_TYPE_MARKET":        1,
		"ORDER_TYPE_LIMIT":         2,
		"ORDER_TYPE_STOP":          3,
		"ORDER_TYPE_STOP_LIMIT":    4,
		"ORDER_TYPE_TRAILING_STOP": 5,
	}
)

func (x OrderType) Enum() *OrderType {
	p := new(OrderType)
	*p = x
	return p
}

func (x OrderType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderType) Descriptor() protoreflect.EnumDescriptor {
	return file_tradesvc_v1_tradesvc_proto_enumTypes[1].Descriptor()
}

func (OrderType) Type() protoreflect.EnumType {
	return &file_tradesvc_v1_tradesvc_proto_enumTypes[1]
}

func (x OrderType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderType.Descriptor instead.
func (OrderType) EnumDescriptor() ([]byte, []int) {
	return file_tradesvc_v1_tradesvc_proto_rawDescGZIP(), []int{1}
}

type OrderClass int32

const (
	OrderClass_ORDER_CLASS_UNSPECIFIED OrderClass = 0
	OrderClass_ORDER_CLASS_BRACKET     OrderClass = 1
	OrderClass_ORDER_CLASS_OTO         OrderClass = 2
	OrderClass_ORDER_CLASS_OCO         OrderClass = 3
	OrderClass_ORDER_CLASS_SIMPLE      OrderClass = 4
)

// Enum value maps for OrderClass.
var (
	OrderClass_name = map[int32]string{
		0: "ORDER_CLASS_UNSPECIFIED",
		1: "ORDER_CLASS_BRACKET",
		2: "ORDER_CLASS_OTO",
		3: "ORDER_CLASS_OCO",
		4: "ORDER_CLASS_SIMPLE",
	}
	OrderClass_value = map[string]int32{
		"ORDER_CLASS_UNSPECIFIED": 0,
		"ORDER_CLASS_BRACKET":     1,
		"ORDER_CLASS_OTO":         2,
		"ORDER_CLASS_OCO":         3,
		"ORDER_CLASS_SIMPLE":      4,
	}
)

func (x OrderClass) Enum() *OrderClass {
	p := new(OrderClass)
	*p = x
	return p
}

func (x OrderClass) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderClass) Descriptor() protoreflect.EnumDescriptor {
	return file_tradesvc_v1_tradesvc_proto_enumTypes[2].Descriptor()
}

func (OrderClass) Type() protoreflect.EnumType {
	return &file_tradesvc_v1_tradesvc_proto_enumTypes[2]
}

func (x OrderClass) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderClass.Descriptor instead.
func (OrderClass) EnumDescriptor() ([]byte, []int) {
	return file_tradesvc_v1_tradesvc_proto_rawDescGZIP(), []int{2}
}

type TimeInForce int32

const (
	TimeInForce_TIME_IN_FORCE_UNSPECIFIED TimeInForce = 0
	TimeInForce_TIME_IN_FORCE_DAY         TimeInForce = 1
	TimeInForce_TIME_IN_FORCE_GTC         TimeInForce = 2
	TimeInForce_TIME_IN_FORCE_OPG         TimeInForce = 3
	TimeInForce_TIME_IN_FORCE_IOC         TimeInForce = 4
	TimeInForce_TIME_IN_FORCE_FOK         TimeInForce = 5
	TimeInForce_TIME_IN_FORCE_GTX         TimeInForce = 6
	TimeInForce_TIME_IN_FORCE_GTD         TimeInForce = 7
	TimeInForce_TIME_IN_FORCE_CLS         TimeInForce = 8
)

// Enum value maps for TimeInForce.
var (
	TimeInForce_name = map[int32]string{
		0: "TIME_IN_FORCE_UNSPECIFIED",
		1: "TIME_IN_FORCE_DAY",
		2: "TIME_IN_FORCE_GTC",
		3: "TIME_IN_FORCE_OPG",
		4: "TIME_IN_FORCE_IOC",
		5: "TIME_IN_FORCE_FOK",
		6: "TIME_IN_FORCE_GTX",
		7: "TIME_IN_FORCE_GTD",
		8: "TIME_IN_FORCE_CLS",
	}
	TimeInForce_value = map[string]int32{
		"TIME_IN_FORCE_UNSPECIFIED": 0,
		"TIME_IN_FORCE_DAY":         1,
		"TIME_IN_FORCE_GTC":         2,
		"TIME_IN_FORCE_OPG":         3,
		"TIME_IN_FORCE_IOC":         4,
		"TIME_IN_FORCE_FOK":         5,
		"TIME_IN_FORCE_GTX":         6,
		"TIME_IN_FORCE_GTD":         7,
		"TIME_IN_FORCE_CLS":         8,
	}
)

func (x TimeInForce) Enum() *TimeInForce {
	p := new(TimeInForce)
	*p = x
	return p
}

func (x TimeInForce) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TimeInForce) Descriptor() protoreflect.EnumDescriptor {
	return file_tradesvc_v1_tradesvc_proto_enumTypes[3].Descriptor()
}

func (TimeInForce) Type() protoreflect.EnumType {
	return &file_tradesvc_v1_tradesvc_proto_enumTypes[3]
}

func (x TimeInForce) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TimeInForce.Descriptor instead.
func (TimeInForce) EnumDescriptor() ([]byte, []int) {
	return file_tradesvc_v1_tradesvc_proto_rawDescGZIP(), []int{3}
}

type PositionIntent int32

const (
	PositionIntent_POSITION_INTENT_UNSPECIFIED   PositionIntent = 0
	PositionIntent_POSITION_INTENT_BUY_TO_OPEN   PositionIntent = 1
	PositionIntent_POSITION_INTENT_BUY_TO_CLOSE  PositionIntent = 2
	PositionIntent_POSITION_INTENT_SELL_TO_OPEN  PositionIntent = 3
	PositionIntent_POSITION_INTENT_SELL_TO_CLOSE PositionIntent = 4
)

// Enum value maps for PositionIntent.
var (
	PositionIntent_name = map[int32]string{
		0: "POSITION_INTENT_UNSPECIFIED",
		1: "POSITION_INTENT_BUY_TO_OPEN",
		2: "POSITION_INTENT_BUY_TO_CLOSE",
		3: "POSITION_INTENT_SELL_TO_OPEN",
		4: "POSITION_INTENT_SELL_TO_CLOSE",
	}
	PositionIntent_value = map[string]int32{
		"POSITION_INTENT_UNSPECIFIED":   0,
		"POSITION_INTENT_BUY_TO_OPEN":   1,
		"POSITION_INTENT_BUY_TO_CLOSE":  2,
		"POSITION_INTENT_SELL_TO_OPEN":  3,
		"POSITION_INTENT_SELL_TO_CLOSE": 4,
	}
)

func (x PositionIntent) Enum() *PositionIntent {
	p := new(PositionIntent)
	*p = x
	return p
}

func (x PositionIntent) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PositionIntent) Descriptor() protoreflect.EnumDescriptor {
	return file_tradesvc_v1_tradesvc_proto_enumTypes[4].Descriptor()
}

func (PositionIntent) Type() protoreflect.EnumType {
	return &file_tradesvc_v1_tradesvc_proto_enumTypes[4]
}

func (x PositionIntent) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PositionIntent.Descriptor instead.
func (PositionIntent) EnumDescriptor() ([]byte, []int) {
	return file_tradesvc_v1_tradesvc_proto_rawDescGZIP(), []int{4}
}

type TakeProfit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LimitPrice float64 `protobuf:"fixed64,1,opt,name=limit_price,json=limitPrice,proto3" json:"limit_price,omitempty"`
}

func (x *TakeProfit) Reset() {
	*x = TakeProfit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradesvc_v1_tradesvc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TakeProfit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TakeProfit) ProtoMessage() {}

func (x *TakeProfit) ProtoReflect() protoreflect.Message {
	mi := &file_tradesvc_v1_tradesvc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TakeProfit.ProtoReflect.Descriptor instead.
func (*TakeProfit) Descriptor() ([]byte, []int) {
	return file_tradesvc_v1_tradesvc_proto_rawDescGZIP(), []int{0}
}

func (x *TakeProfit) GetLimitPrice() float64 {
	if x != nil {
		return x.LimitPrice
	}
	return 0
}

type StopLoss struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit float64 `protobuf:"fixed64,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Stop  float64 `protobuf:"fixed64,2,opt,name=stop,proto3" json:"stop,omitempty"`
}

func (x *StopLoss) Reset() {
	*x = StopLoss{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradesvc_v1_tradesvc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopLoss) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopLoss) ProtoMessage() {}

func (x *StopLoss) ProtoReflect() protoreflect.Message {
	mi := &file_tradesvc_v1_tradesvc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopLoss.ProtoReflect.Descriptor instead.
func (*StopLoss) Descriptor() ([]byte, []int) {
	return file_tradesvc_v1_tradesvc_proto_rawDescGZIP(), []int{1}
}

func (x *StopLoss) GetLimit() float64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *StopLoss) GetStop() float64 {
	if x != nil {
		return x.Stop
	}
	return 0
}

type Trade struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol         string         `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Qty            float64        `protobuf:"fixed64,2,opt,name=qty,proto3" json:"qty,omitempty"`
	Notional       float64        `protobuf:"fixed64,3,opt,name=notional,proto3" json:"notional,omitempty"`
	Side           Side           `protobuf:"varint,4,opt,name=side,proto3,enum=tradesvc.v1.Side" json:"side,omitempty"`
	OrderType      OrderType      `protobuf:"varint,5,opt,name=order_type,json=orderType,proto3,enum=tradesvc.v1.OrderType" json:"order_type,omitempty"`
	Tif            TimeInForce    `protobuf:"varint,6,opt,name=tif,proto3,enum=tradesvc.v1.TimeInForce" json:"tif,omitempty"`
	LimitPrice     float64        `protobuf:"fixed64,7,opt,name=limit_price,json=limitPrice,proto3" json:"limit_price,omitempty"`
	ExtendedHours  bool           `protobuf:"varint,8,opt,name=extended_hours,json=extendedHours,proto3" json:"extended_hours,omitempty"`
	StopPrice      float64        `protobuf:"fixed64,9,opt,name=stop_price,json=stopPrice,proto3" json:"stop_price,omitempty"`
	Class          OrderClass     `protobuf:"varint,11,opt,name=class,proto3,enum=tradesvc.v1.OrderClass" json:"class,omitempty"`
	TakeProfit     *TakeProfit    `protobuf:"bytes,12,opt,name=take_profit,json=takeProfit,proto3" json:"take_profit,omitempty"`
	StopLoss       *StopLoss      `protobuf:"bytes,13,opt,name=stop_loss,json=stopLoss,proto3" json:"stop_loss,omitempty"`
	TrailPrice     float64        `protobuf:"fixed64,14,opt,name=trail_price,json=trailPrice,proto3" json:"trail_price,omitempty"`
	TrailPercent   float64        `protobuf:"fixed64,15,opt,name=trail_percent,json=trailPercent,proto3" json:"trail_percent,omitempty"`
	PositionIntent PositionIntent `protobuf:"varint,16,opt,name=position_intent,json=positionIntent,proto3,enum=tradesvc.v1.PositionIntent" json:"position_intent,omitempty"`
}

func (x *Trade) Reset() {
	*x = Trade{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tradesvc_v1_tradesvc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Trade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Trade) ProtoMessage() {}

func (x *Trade) ProtoReflect() protoreflect.Message {
	mi := &file_tradesvc_v1_tradesvc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Trade.ProtoReflect.Descriptor instead.
func (*Trade) Descriptor() ([]byte, []int) {
	return file_tradesvc_v1_tradesvc_proto_rawDescGZIP(), []int{2}
}

func (x *Trade) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Trade) GetQty() float64 {
	if x != nil {
		return x.Qty
	}
	return 0
}

func (x *Trade) GetNotional() float64 {
	if x != nil {
		return x.Notional
	}
	return 0
}

func (x *Trade) GetSide() Side {
	if x != nil {
		return x.Side
	}
	return Side_SIDE_UNSPECIFIED
}

func (x *Trade) GetOrderType() OrderType {
	if x != nil {
		return x.OrderType
	}
	return OrderType_ORDER_TYPE_UNSPECIFIED
}

func (x *Trade) GetTif() TimeInForce {
	if x != nil {
		return x.Tif
	}
	return TimeInForce_TIME_IN_FORCE_UNSPECIFIED
}

func (x *Trade) GetLimitPrice() float64 {
	if x != nil {
		return x.LimitPrice
	}
	return 0
}

func (x *Trade) GetExtendedHours() bool {
	if x != nil {
		return x.ExtendedHours
	}
	return false
}

func (x *Trade) GetStopPrice() float64 {
	if x != nil {
		return x.StopPrice
	}
	return 0
}

func (x *Trade) GetClass() OrderClass {
	if x != nil {
		return x.Class
	}
	return OrderClass_ORDER_CLASS_UNSPECIFIED
}

func (x *Trade) GetTakeProfit() *TakeProfit {
	if x != nil {
		return x.TakeProfit
	}
	return nil
}

func (x *Trade) GetStopLoss() *StopLoss {
	if x != nil {
		return x.StopLoss
	}
	return nil
}

func (x *Trade) GetTrailPrice() float64 {
	if x != nil {
		return x.TrailPrice
	}
	return 0
}

func (x *Trade) GetTrailPercent() float64 {
	if x != nil {
		return x.TrailPercent
	}
	return 0
}

func (x *Trade) GetPositionIntent() PositionIntent {
	if x != nil {
		return x.PositionIntent
	}
	return PositionIntent_POSITION_INTENT_UNSPECIFIED
}

var File_tradesvc_v1_tradesvc_proto protoreflect.FileDescriptor

var file_tradesvc_v1_tradesvc_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x74, 0x72, 0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72,
	0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x74, 0x72,
	0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x22, 0x2d, 0x0a, 0x0a, 0x54, 0x61, 0x6b,
	0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x22, 0x34, 0x0a, 0x08, 0x53, 0x74, 0x6f, 0x70,
	0x4c, 0x6f, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x74,
	0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x73, 0x74, 0x6f, 0x70, 0x22, 0xe7,
	0x04, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x12, 0x10, 0x0a, 0x03, 0x71, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x71,
	0x74, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x6f, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6e, 0x6f, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x12, 0x25,
	0x0a, 0x04, 0x73, 0x69, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x74,
	0x72, 0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x64, 0x65, 0x52,
	0x04, 0x73, 0x69, 0x64, 0x65, 0x12, 0x35, 0x0a, 0x0a, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x74, 0x72, 0x61, 0x64,
	0x65, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x09, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2a, 0x0a, 0x03,
	0x74, 0x69, 0x66, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x74, 0x72, 0x61, 0x64,
	0x65, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x46, 0x6f,
	0x72, 0x63, 0x65, 0x52, 0x03, 0x74, 0x69, 0x66, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x64, 0x65, 0x64, 0x5f, 0x68, 0x6f, 0x75, 0x72, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0d, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x48, 0x6f, 0x75, 0x72, 0x73,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x6f, 0x70, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x70, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12,
	0x2d, 0x0a, 0x05, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17,
	0x2e, 0x74, 0x72, 0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x52, 0x05, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x38,
	0x0a, 0x0b, 0x74, 0x61, 0x6b, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x74, 0x72, 0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x61, 0x6b, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x52, 0x0a, 0x74, 0x61,
	0x6b, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x12, 0x32, 0x0a, 0x09, 0x73, 0x74, 0x6f, 0x70,
	0x5f, 0x6c, 0x6f, 0x73, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x72,
	0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x4c, 0x6f,
	0x73, 0x73, 0x52, 0x08, 0x73, 0x74, 0x6f, 0x70, 0x4c, 0x6f, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x0b,
	0x74, 0x72, 0x61, 0x69, 0x6c, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0a, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a,
	0x0d, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x50, 0x65, 0x72, 0x63, 0x65,
	0x6e, 0x74, 0x12, 0x44, 0x0a, 0x0f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x74, 0x72,
	0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x0e, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2a, 0x39, 0x0a, 0x04, 0x53, 0x69, 0x64, 0x65,
	0x12, 0x14, 0x0a, 0x10, 0x53, 0x49, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x49, 0x44, 0x45, 0x5f, 0x42,
	0x55, 0x59, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x49, 0x44, 0x45, 0x5f, 0x53, 0x45, 0x4c,
	0x4c, 0x10, 0x02, 0x2a, 0xa2, 0x01, 0x0a, 0x09, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1a, 0x0a, 0x16, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a,
	0x11, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x41, 0x52, 0x4b,
	0x45, 0x54, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x4f, 0x52,
	0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x4f, 0x50, 0x10, 0x03, 0x12,
	0x19, 0x0a, 0x15, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54,
	0x4f, 0x50, 0x5f, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x10, 0x04, 0x12, 0x1c, 0x0a, 0x18, 0x4f, 0x52,
	0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x52, 0x41, 0x49, 0x4c, 0x49, 0x4e,
	0x47, 0x5f, 0x53, 0x54, 0x4f, 0x50, 0x10, 0x05, 0x2a, 0x84, 0x01, 0x0a, 0x0a, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x1b, 0x0a, 0x17, 0x4f, 0x52, 0x44, 0x45, 0x52,
	0x5f, 0x43, 0x4c, 0x41, 0x53, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x43, 0x4c,
	0x41, 0x53, 0x53, 0x5f, 0x42, 0x52, 0x41, 0x43, 0x4b, 0x45, 0x54, 0x10, 0x01, 0x12, 0x13, 0x0a,
	0x0f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x43, 0x4c, 0x41, 0x53, 0x53, 0x5f, 0x4f, 0x54, 0x4f,
	0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x43, 0x4c, 0x41, 0x53,
	0x53, 0x5f, 0x4f, 0x43, 0x4f, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x4f, 0x52, 0x44, 0x45, 0x52,
	0x5f, 0x43, 0x4c, 0x41, 0x53, 0x53, 0x5f, 0x53, 0x49, 0x4d, 0x50, 0x4c, 0x45, 0x10, 0x04, 0x2a,
	0xe4, 0x01, 0x0a, 0x0b, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x12,
	0x1d, 0x0a, 0x19, 0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15,
	0x0a, 0x11, 0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x5f,
	0x44, 0x41, 0x59, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e,
	0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x5f, 0x47, 0x54, 0x43, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11,
	0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x5f, 0x4f, 0x50,
	0x47, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46,
	0x4f, 0x52, 0x43, 0x45, 0x5f, 0x49, 0x4f, 0x43, 0x10, 0x04, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x49,
	0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x5f, 0x46, 0x4f, 0x4b, 0x10,
	0x05, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52,
	0x43, 0x45, 0x5f, 0x47, 0x54, 0x58, 0x10, 0x06, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x49, 0x4d, 0x45,
	0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45, 0x5f, 0x47, 0x54, 0x44, 0x10, 0x07, 0x12,
	0x15, 0x0a, 0x11, 0x54, 0x49, 0x4d, 0x45, 0x5f, 0x49, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x43, 0x45,
	0x5f, 0x43, 0x4c, 0x53, 0x10, 0x08, 0x2a, 0xb9, 0x01, 0x0a, 0x0e, 0x50, 0x6f, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x1b, 0x50, 0x4f, 0x53,
	0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1f, 0x0a, 0x1b, 0x50, 0x4f,
	0x53, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x42, 0x55,
	0x59, 0x5f, 0x54, 0x4f, 0x5f, 0x4f, 0x50, 0x45, 0x4e, 0x10, 0x01, 0x12, 0x20, 0x0a, 0x1c, 0x50,
	0x4f, 0x53, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x42,
	0x55, 0x59, 0x5f, 0x54, 0x4f, 0x5f, 0x43, 0x4c, 0x4f, 0x53, 0x45, 0x10, 0x02, 0x12, 0x20, 0x0a,
	0x1c, 0x50, 0x4f, 0x53, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54,
	0x5f, 0x53, 0x45, 0x4c, 0x4c, 0x5f, 0x54, 0x4f, 0x5f, 0x4f, 0x50, 0x45, 0x4e, 0x10, 0x03, 0x12,
	0x21, 0x0a, 0x1d, 0x50, 0x4f, 0x53, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x4e, 0x54, 0x45,
	0x4e, 0x54, 0x5f, 0x53, 0x45, 0x4c, 0x4c, 0x5f, 0x54, 0x4f, 0x5f, 0x43, 0x4c, 0x4f, 0x53, 0x45,
	0x10, 0x04, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x41, 0x6e, 0x74, 0x68, 0x6f, 0x6e, 0x79, 0x48, 0x65, 0x77, 0x69, 0x6e, 0x73, 0x2f, 0x66,
	0x61, 0x6c, 0x70, 0x61, 0x63, 0x61, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x74, 0x72,
	0x61, 0x64, 0x65, 0x73, 0x76, 0x63, 0x73, 0x76, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x74, 0x72, 0x61,
	0x64, 0x65, 0x73, 0x76, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tradesvc_v1_tradesvc_proto_rawDescOnce sync.Once
	file_tradesvc_v1_tradesvc_proto_rawDescData = file_tradesvc_v1_tradesvc_proto_rawDesc
)

func file_tradesvc_v1_tradesvc_proto_rawDescGZIP() []byte {
	file_tradesvc_v1_tradesvc_proto_rawDescOnce.Do(func() {
		file_tradesvc_v1_tradesvc_proto_rawDescData = protoimpl.X.CompressGZIP(file_tradesvc_v1_tradesvc_proto_rawDescData)
	})
	return file_tradesvc_v1_tradesvc_proto_rawDescData
}

var file_tradesvc_v1_tradesvc_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_tradesvc_v1_tradesvc_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_tradesvc_v1_tradesvc_proto_goTypes = []interface{}{
	(Side)(0),           // 0: tradesvc.v1.Side
	(OrderType)(0),      // 1: tradesvc.v1.OrderType
	(OrderClass)(0),     // 2: tradesvc.v1.OrderClass
	(TimeInForce)(0),    // 3: tradesvc.v1.TimeInForce
	(PositionIntent)(0), // 4: tradesvc.v1.PositionIntent
	(*TakeProfit)(nil),  // 5: tradesvc.v1.TakeProfit
	(*StopLoss)(nil),    // 6: tradesvc.v1.StopLoss
	(*Trade)(nil),       // 7: tradesvc.v1.Trade
}
var file_tradesvc_v1_tradesvc_proto_depIdxs = []int32{
	0, // 0: tradesvc.v1.Trade.side:type_name -> tradesvc.v1.Side
	1, // 1: tradesvc.v1.Trade.order_type:type_name -> tradesvc.v1.OrderType
	3, // 2: tradesvc.v1.Trade.tif:type_name -> tradesvc.v1.TimeInForce
	2, // 3: tradesvc.v1.Trade.class:type_name -> tradesvc.v1.OrderClass
	5, // 4: tradesvc.v1.Trade.take_profit:type_name -> tradesvc.v1.TakeProfit
	6, // 5: tradesvc.v1.Trade.stop_loss:type_name -> tradesvc.v1.StopLoss
	4, // 6: tradesvc.v1.Trade.position_intent:type_name -> tradesvc.v1.PositionIntent
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_tradesvc_v1_tradesvc_proto_init() }
func file_tradesvc_v1_tradesvc_proto_init() {
	if File_tradesvc_v1_tradesvc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tradesvc_v1_tradesvc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TakeProfit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tradesvc_v1_tradesvc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopLoss); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tradesvc_v1_tradesvc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Trade); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tradesvc_v1_tradesvc_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tradesvc_v1_tradesvc_proto_goTypes,
		DependencyIndexes: file_tradesvc_v1_tradesvc_proto_depIdxs,
		EnumInfos:         file_tradesvc_v1_tradesvc_proto_enumTypes,
		MessageInfos:      file_tradesvc_v1_tradesvc_proto_msgTypes,
	}.Build()
	File_tradesvc_v1_tradesvc_proto = out.File
	file_tradesvc_v1_tradesvc_proto_rawDesc = nil
	file_tradesvc_v1_tradesvc_proto_goTypes = nil
	file_tradesvc_v1_tradesvc_proto_depIdxs = nil
}
