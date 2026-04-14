/*
 * Copyright 2024 Tech4Works
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package publisher

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	snsTypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"

	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/publisher"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"github.com/tech4works/gopen-gateway/internal/infra/telemetry"
)

type client struct {
	sqs *sqs.Client
	sns *sns.Client
}

type messageAttribute struct {
	dataType string
	value    string
	binary   []byte
}

const (
	dataTypeString = "String"
	dataTypeNumber = "Number"
	dataTypeBinary = "Binary"
)

func NewClient(sqsClient *sqs.Client, snsClient *sns.Client) app.PublisherClient {
	return client{
		sqs: sqsClient,
		sns: snsClient,
	}
}

func (c client) Publish(
	ctx context.Context,
	parent *vo.EndpointRequest,
	request *vo.PublisherBackendRequest,
) (*publisher.Response, error) {
	ctx, span := telemetry.Tracer().Start(ctx, "messaging.publish")
	defer span.End()

	c.fillSpanAttributes(span, request)

	var resp *publisher.Response
	var err error

	switch request.Broker() {
	case enum.BackendBrokerAwsSqs:
		resp, err = c.publishSQS(ctx, parent, request)
	case enum.BackendBrokerAwsSns:
		resp, err = c.publishSNS(ctx, parent, request)
	default:
		err = app.NewErrBackendBrokerNotImplemented(request.Broker().String())
	}

	if checker.NonNil(err) {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return resp, nil
}

func (c client) publishSQS(
	ctx context.Context,
	parent *vo.EndpointRequest,
	request *vo.PublisherBackendRequest,
) (*publisher.Response, error) {
	if checker.IsNil(c.sqs) {
		return nil, app.NewErrBackendBrokerNotConfigured(enum.BackendBrokerAwsSqs.String())
	}

	inp, err := c.buildSQSPublishInput(parent, request)
	if checker.NonNil(err) {
		return nil, errors.Inheritf(err, "publish client: op=build-sqs-publish-input")
	}

	out, err := c.sqs.SendMessage(ctx, inp)
	if checker.NonNil(err) {
		return nil, err
	}

	return c.buildResponse(request, out.MessageId, out.SequenceNumber), nil
}

func (c client) publishSNS(
	ctx context.Context,
	parent *vo.EndpointRequest,
	request *vo.PublisherBackendRequest,
) (*publisher.Response, error) {
	if checker.IsNil(c.sns) {
		return nil, app.NewErrBackendBrokerNotConfigured(enum.BackendBrokerAwsSns.String())
	}

	inp, err := c.buildSNSPublishInput(parent, request)
	if checker.NonNil(err) {
		return nil, errors.Inheritf(err, "publish client: op=build-sns-publish-input")
	}

	out, err := c.sns.Publish(ctx, inp)
	if checker.NonNil(err) {
		return nil, err
	}

	return c.buildResponse(request, out.MessageId, out.SequenceNumber), nil
}

func (c client) buildResponse(
	request *vo.PublisherBackendRequest,
	messageID *string,
	sequenceNumber *string,
) *publisher.Response {
	response := &publisher.Response{
		OK: checker.IsNotNilOrEmpty(messageID),
		Body: &publisher.Body{
			Path:     request.Path(),
			Provider: request.Broker().String(),
		},
	}

	if checker.IsNotNilOrEmpty(messageID) {
		response.Body.MessageID = *messageID
	}
	if checker.IsNotNilOrEmpty(sequenceNumber) {
		response.Body.SequentialNumber = *sequenceNumber
	}

	return response
}

func (c client) buildSQSPublishInput(
	parent *vo.EndpointRequest,
	request *vo.PublisherBackendRequest,
) (*sqs.SendMessageInput, error) {
	body, err := c.parseBodyToPointerString(request)
	if checker.NonNil(err) {
		return nil, err
	}

	return &sqs.SendMessageInput{
		QueueUrl:               converter.ToPointer(request.Path()),
		MessageBody:            body,
		DelaySeconds:           int32(request.Delay().Time().Seconds()),
		MessageAttributes:      c.toSQSMessageAttributes(c.buildMessageAttributes(parent, request)),
		MessageDeduplicationId: request.DeduplicationID(),
		MessageGroupId:         request.GroupID(),
	}, nil
}

func (c client) buildSNSPublishInput(
	parent *vo.EndpointRequest,
	request *vo.PublisherBackendRequest,
) (*sns.PublishInput, error) {
	body, err := c.parseBodyToPointerString(request)
	if checker.NonNil(err) {
		return nil, err
	}

	return &sns.PublishInput{
		TopicArn:               converter.ToPointer(request.Path()),
		Message:                body,
		MessageAttributes:      c.toSNSMessageAttributes(c.buildMessageAttributes(parent, request)),
		MessageDeduplicationId: request.DeduplicationID(),
		MessageGroupId:         request.GroupID(),
	}, nil
}

func (c client) buildMessageAttributes(
	parent *vo.EndpointRequest,
	request *vo.PublisherBackendRequest,
) map[string]messageAttribute {
	messageAttributes := map[string]messageAttribute{}

	for key, attribute := range request.Attributes() {
		switch attribute.Type() {
		case enum.AttributeValueTypeString:
			messageAttributes[key] = messageAttribute{
				dataType: dataTypeString,
				value:    attribute.Value(),
			}
		case enum.AttributeValueTypeNumber:
			messageAttributes[key] = messageAttribute{
				dataType: dataTypeNumber,
				value:    attribute.Value(),
			}
		case enum.AttributeValueTypeBinary:
			messageAttributes[key] = messageAttribute{
				dataType: dataTypeBinary,
				binary:   converter.ToBytes(attribute.Value()),
			}
		}
	}

	c.appendDefaultMessageAttributes(messageAttributes, parent, request)

	return messageAttributes
}

func (c client) appendDefaultMessageAttributes(
	messageAttributes map[string]messageAttribute,
	parent *vo.EndpointRequest,
	request *vo.PublisherBackendRequest,
) {
	messageAttributes[app.XGopenRequestID] = messageAttribute{
		dataType: dataTypeString,
		value:    parent.ID(),
	}
	messageAttributes[app.XForwardedFor] = messageAttribute{
		dataType: dataTypeString,
		value:    parent.ClientIP(),
	}
	messageAttributes[app.XGopenDegraded] = messageAttribute{
		dataType: dataTypeString,
		value:    converter.ToString(request.Degraded()),
	}
	messageAttributes[app.XGopenGroupIDDegraded] = messageAttribute{
		dataType: dataTypeString,
		value:    converter.ToString(request.GroupIDDegraded()),
	}
	messageAttributes[app.XGopenDeduplicationIDDegraded] = messageAttribute{
		dataType: dataTypeString,
		value:    converter.ToString(request.DeduplicationIDDegraded()),
	}
	messageAttributes[app.XGopenAttributeDegraded] = messageAttribute{
		dataType: dataTypeString,
		value:    converter.ToString(request.AttributesDegraded()),
	}
	messageAttributes[app.XGopenBodyDegraded] = messageAttribute{
		dataType: dataTypeString,
		value:    converter.ToString(request.BodyDegraded()),
	}
}

func (c client) toSQSMessageAttributes(
	attributes map[string]messageAttribute,
) map[string]sqsTypes.MessageAttributeValue {
	out := make(map[string]sqsTypes.MessageAttributeValue, len(attributes))

	for key, attribute := range attributes {
		value := sqsTypes.MessageAttributeValue{
			DataType: converter.ToPointer(attribute.dataType),
		}

		if attribute.dataType == "Binary" {
			value.BinaryValue = attribute.binary
		} else {
			value.StringValue = converter.ToPointer(attribute.value)
		}

		out[key] = value
	}

	return out
}

func (c client) toSNSMessageAttributes(
	attributes map[string]messageAttribute,
) map[string]snsTypes.MessageAttributeValue {
	out := make(map[string]snsTypes.MessageAttributeValue, len(attributes))

	for key, attribute := range attributes {
		value := snsTypes.MessageAttributeValue{
			DataType: converter.ToPointer(attribute.dataType),
		}

		if checker.EqualsIgnoreCase(attribute.dataType, "Binary") {
			value.BinaryValue = attribute.binary
		} else {
			value.StringValue = converter.ToPointer(attribute.value)
		}

		out[key] = value
	}

	return out
}

func (c client) parseBodyToPointerString(request *vo.PublisherBackendRequest) (*string, error) {
	if !request.HasBody() {
		return nil, nil
	}

	compactString, err := request.Body().CompactString()
	if checker.NonNil(err) {
		return nil, errors.Inheritf(err, "publish client: op=parse-body-compact-string")
	}

	return converter.ToPointer(compactString), nil
}

func (c client) fillSpanAttributes(span trace.Span, request *vo.PublisherBackendRequest) {
	span.SetAttributes(
		attribute.String("messaging.broker", request.Broker().String()),
		attribute.String("messaging.path", request.Path()),
	)

	if checker.IsNil(request.GroupID()) {
		span.SetAttributes(attribute.String("messaging.group_id", "<nil>"))
	} else {
		span.SetAttributes(attribute.String("messaging.group_id", *request.GroupID()))
	}

	if checker.IsNil(request.DeduplicationID()) {
		span.SetAttributes(attribute.String("messaging.deduplication_id", "<nil>"))
	} else {
		span.SetAttributes(attribute.String("messaging.deduplication_id", *request.DeduplicationID()))
	}

	if request.HasBody() {
		bodyCompactString, err := request.Body().CompactString()
		if checker.NonNil(err) {
			span.SetAttributes(attribute.String("messaging.body", err.Error()))
		} else {
			span.SetAttributes(attribute.String("messaging.body", bodyCompactString))
		}
	} else {
		span.SetAttributes(attribute.String("messaging.body", "<nil>"))
	}
}
