/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package notify

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Channel represents a message transmission channel in TKE.
type Channel struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// Spec defines the desired channel.
	// +optional
	Spec ChannelSpec
	// +optional
	Status ChannelStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChannelList is the whole list of all channels which owned by a tenant.
type ChannelList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of channels.
	Items []Channel
}

// FinalizerName is the name identifying a finalizer during channel lifecycle.
type FinalizerName string

const (
	// ChannelFinalize is an internal finalizer values to Channel.
	ChannelFinalize FinalizerName = "channel"
)

// ChannelSpec is a description of a channel.
type ChannelSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers  []FinalizerName
	TenantID    string
	DisplayName string
	// +optional
	TencentCloudSMS *ChannelTencentCloudSMS
	// +optional
	Wechat *ChannelWechat
	// +optional
	SMTP *ChannelSMTP
}

// ChannelStatus represents information about the status of a cluster.
type ChannelStatus struct {
	// +optional
	Phase ChannelPhase `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase,casttype=ChannelPhase"`
}

// ChannelPhase defines the phase of channel constructor.
type ChannelPhase string

const (
	// ChannelActived is the normal running phase.
	ChannelActived ChannelPhase = "Actived"
	// ChannelTerminating means the channel is undergoing graceful termination.
	ChannelTerminating ChannelPhase = "Terminating"
)

// ChannelTencentCloudSMS indicates the channel configuration for sending
// messages using Tencent Cloud SMS Gateway.
// See: https://cloud.tencent.com/document/product/382/5976
type ChannelTencentCloudSMS struct {
	AppKey   string
	SdkAppID string
	// +optional
	Extend string
}

// ChannelWechat indicates a channel configuration for sending template
// notifications using WeChat.
type ChannelWechat struct {
	// AppID indicates the unique credentials of the third-party user.
	// See https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183
	AppID     string
	AppSecret string
}

// ChannelSMTP indicates a channel configuration for sending email notifications
// using the SMTP server.
type ChannelSMTP struct {
	SMTPHost string
	SMTPPort int32
	TLS      bool
	Email    string
	Password string
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Template indicates the template used to send notifications under this channel.
type Template struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// Spec defines the desired template.
	// +optional
	Spec TemplateSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplateList is the whole list of all template which owned by a channel.
type TemplateList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of templates.
	Items []Template
}

// TemplateSpec is a description of a template.
type TemplateSpec struct {
	TenantID    string
	DisplayName string
	// +optional
	Keys []string
	// +optional
	TencentCloudSMS *TemplateTencentCloudSMS
	// +optional
	Wechat *TemplateWechat
	// +optional
	Text *TemplateText
}

// TemplateTencentCloudSMS indicates the template used when sending text
// messages using Tencent Cloud SMS Gateway.
// The template must be approved.
type TemplateTencentCloudSMS struct {
	TemplateID string
	// +optional
	Sign string
	// +optional
	Body string
}

// TemplateWechat indicates the template when sending a text message using the
// WeChat public account.
// The template must be approved and registered.
type TemplateWechat struct {
	// TemplateID indicates the template id of the template message notification.
	// See https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183
	TemplateID string
	// URL indicates the web address of the user who clicked the notification in WeChat.
	// +optional
	URL string
	// MiniProgramAppID indicates the unique identification number of the WeChat applet
	// that the user clicked on the notification in WeChat.
	// +optional
	MiniProgramAppID string
	// +optional
	MiniProgramPagePath string
	// +optional
	Body string
}

// TemplateText indicates the template used to send notifications using other
// non-templated channels.
type TemplateText struct {
	Body string
	// +optional
	Header string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Receiver indicates a message notification recipient, usually representing a
// user in the user system or a webhook service address.
type Receiver struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// Spec defines the desired receiver.
	// +optional
	Spec ReceiverSpec
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReceiverList is the whole list of all receiver which owned by a tenant.
type ReceiverList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of receivers.
	Items []Receiver
}

// ReceiverChannel is the name identifying various channel in a receiver.
type ReceiverChannel string

const (
	// ReceiverChannelMobile represents the mobile of receiver.
	ReceiverChannelMobile ReceiverChannel = "mobile"
	// ReceiverChannelEmail represents the email address of receiver.
	ReceiverChannelEmail ReceiverChannel = "email"
	// ReceiverChannelWechatOpenID represents the openid for wechat of receiver.
	ReceiverChannelWechatOpenID ReceiverChannel = "wechat_openid"
)

// ReceiverSpec is a description of a receiver.
type ReceiverSpec struct {
	TenantID    string
	DisplayName string
	// +optional
	Username string
	// Identities represents the characteristics of the message recipient.
	// The hash table key represents the message delivery channel id, and the value represents
	// the user identification number in the channel.
	// For example, if it is a short message sending channel, then the value is the user's
	// mobile phone number; if it is a mail sending channel, then the value is the user's
	// email address.
	// +optional
	Identities map[ReceiverChannel]string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReceiverGroup indicates multiple message recipients.
type ReceiverGroup struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// Spec defines the desired receiver group.
	// +optional
	Spec ReceiverGroupSpec
}

// ReceiverGroupSpec is a description of a receiver group.
type ReceiverGroupSpec struct {
	TenantID    string
	DisplayName string
	// +optional
	Receivers []string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReceiverGroupList is the whole list of all receiver which owned by a tenant.
type ReceiverGroupList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of receiver groups.
	Items []ReceiverGroup
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MessageRequest represents a message sending request, which may include
// multiple recipients and multiple receiving groups.
type MessageRequest struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// Spec defines the desired message.
	// +optional
	Spec MessageRequestSpec
	// +optional
	Status MessageRequestStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MessageRequestList is the whole list of all message which owned by a tenant.
type MessageRequestList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of message requests.
	Items []MessageRequest
}

// MessageRequestSpec is a description of a message request.
type MessageRequestSpec struct {
	TenantID     string
	TemplateName string
	// +optional
	Receivers []string
	// +optional
	ReceiverGroups []string
	// +optional
	Variables map[string]string
}

// MessageRequestStatus represents information about the status of a message request.
type MessageRequestStatus struct {
	// +optional
	Phase MessageRequestPhase
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// A human readable message indicating details about the transition.
	// +optional
	Errors map[string]string
}

// MessageRequestPhase indicates the status of message request.
type MessageRequestPhase string

// These are valid status of message request.
const (
	// MessageRequestPending indicates that the message request has been declared, when the message
	// has not actually been sent.
	MessageRequestPending MessageRequestPhase = "Pending"
	// MessageRequestSending indicated the message is sending.
	MessageRequestSending MessageRequestPhase = "Sending"
	// MessageRequestSent indicates the message has been sent.
	MessageRequestSent MessageRequestPhase = "Sent"
	// MessageRequestFailed indicates that the message failed to be sent.
	MessageRequestFailed MessageRequestPhase = "Failed"
	// MessageRequestPartialFailure indicates that the partial failure to sent.
	MessageRequestPartialFailure MessageRequestPhase = "PartialFailure"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Message indicates a message in the notification system.
type Message struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// Spec defines the desired message.
	// +optional
	Spec MessageSpec
	// +optional
	Status MessageStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MessageList is the whole list of all message which owned by a tenant.
type MessageList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of messages.
	Items []Message
}

// MessageSpec is a description of a message.
type MessageSpec struct {
	TenantID        string
	ReceiverName    string
	ReceiverChannel ReceiverChannel
	Identity        string
	// +optional
	Username string
	// +optional
	Header string
	// +optional
	Body string
	// +optional
	ChannelMessageID string
}

// MessageStatus represents information about the status of a message.
type MessageStatus struct {
	// +optional
	Phase MessagePhase
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
}

// MessagePhase indicates the status of message.
type MessagePhase string

// These are valid status of message.
const (
	// MessageUnread indicates that the message has not been read by the receiver.
	MessageUnread MessagePhase = "Unread"
	// MessageRead indicates that the recipient has read the message.
	MessageRead MessagePhase = "Read"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMap holds configuration data for tke to consume.
type ConfigMap struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string

	// BinaryData contains the binary data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// BinaryData can contain byte sequences that are not in the UTF-8 range.
	// The keys stored in BinaryData must not overlap with the ones in
	// the Data field, this is enforced during validation process.
	// +optional
	BinaryData map[string][]byte
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapList is a resource containing a list of ConfigMap objects.
type ConfigMapList struct {
	metav1.TypeMeta

	// +optional
	metav1.ListMeta

	// Items is the list of ConfigMaps.
	Items []ConfigMap
}
