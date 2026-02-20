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
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
	"github.com/tech4works/gopen-gateway/internal/app/model/publisher"
	"github.com/tech4works/gopen-gateway/internal/domain/model/enum"
	"github.com/tech4works/gopen-gateway/internal/domain/model/vo"
	"go.elastic.co/apm/v2"
)

type client struct {
	sqs *sqs.Client
	sns *sns.Client
}

func NewClient(sqs *sqs.Client, sns *sns.Client) app.PublisherClient {
	return client{
		sqs: sqs,
		sns: sns,
	}
}

func (c client) Publish(ctx context.Context, request *vo.PublisherBackendRequest) (*publisher.Response, error) {
	span, ctx := apm.StartSpan(ctx, "messaging.publish", "publisher")
	defer span.End()

	span.Context.SetLabel("broker", request.Broker())
	span.Context.SetLabel("url", request.Path())
	span.Context.SetLabel("message", request.Body())

	switch request.Broker() {
	case enum.BackendBrokerAwsSqs:
		return c.publishSQS(ctx, request)
	case enum.BackendBrokerAwsSns:
		return c.publishSNS(ctx, request)
	default:
		return nil, errors.Newf("Broker %s not supported", request.Broker())
	}
}

func (c client) publishSQS(ctx context.Context, request *vo.PublisherBackendRequest) (*publisher.Response, error) {
	if checker.IsNil(c.sqs) {
		return nil, errors.New("SQS client not configuration. Please check your configuration.")
	}

	messageAttributes := make(map[string]sqstypes.MessageAttributeValue, len(request.Attributes()))
	for key, attribute := range request.Attributes() {
		messageAttributes[key] = sqstypes.MessageAttributeValue{
			DataType:    converter.ToPointer(attribute.DataType()),
			StringValue: converter.ToPointer(attribute.Value()),
		}
	}

	out, err := c.sqs.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:            converter.ToPointer(request.Body()),
		QueueUrl:               converter.ToPointer(request.Path()),
		DelaySeconds:           int32(request.Delay().Time().Seconds()),
		MessageAttributes:      messageAttributes,
		MessageDeduplicationId: request.DeduplicationID(),
		MessageGroupId:         request.GroupID(),
	})
	if checker.NonNil(err) {
		return nil, err
	}

	return &publisher.Response{
		OK: checker.IsNotNilOrEmpty(out.MessageId),
		Body: &publisher.Body{
			Path:             request.Path(),
			Provider:         request.Broker().String(),
			MessageID:        *out.MessageId,
			SequentialNumber: *out.SequenceNumber,
		},
	}, nil
}

func (c client) publishSNS(ctx context.Context, request *vo.PublisherBackendRequest) (*publisher.Response, error) {
	if checker.IsNil(c.sns) {
		return nil, errors.New("SNS client not configuration. Please check your configuration.")
	}

	messageAttributes := make(map[string]snstypes.MessageAttributeValue, len(request.Attributes()))
	for key, attribute := range request.Attributes() {
		messageAttributes[key] = snstypes.MessageAttributeValue{
			DataType:    converter.ToPointer(attribute.DataType()),
			StringValue: converter.ToPointer(attribute.Value()),
		}
	}

	out, err := c.sns.Publish(ctx, &sns.PublishInput{
		TopicArn:               converter.ToPointer(request.Path()),
		Message:                converter.ToPointer(request.Body()),
		MessageDeduplicationId: request.DeduplicationID(),
		MessageGroupId:         request.GroupID(),
		MessageAttributes:      messageAttributes,
	})
	if checker.NonNil(err) {
		return nil, err
	}

	return &publisher.Response{
		OK: checker.IsNotNilOrEmpty(out.MessageId),
		Body: &publisher.Body{
			Path:             request.Path(),
			Provider:         request.Broker().String(),
			MessageID:        *out.MessageId,
			SequentialNumber: *out.SequenceNumber,
		},
	}, nil
}
