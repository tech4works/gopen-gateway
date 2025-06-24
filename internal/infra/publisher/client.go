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
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/tech4works/checker"
	"github.com/tech4works/converter"
	"github.com/tech4works/errors"
	"github.com/tech4works/gopen-gateway/internal/app"
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

func (c client) Publish(ctx context.Context, publisher *vo.Publisher, message *vo.Message) error {
	span, ctx := apm.StartSpan(ctx, "messaging.publish", "publisher")
	defer span.End()

	span.Context.SetLabel("provider", publisher.Provider())
	span.Context.SetLabel("url", publisher.Reference())
	span.Context.SetLabel("message", message.Body())

	switch publisher.Provider() {
	case enum.PublisherProviderAwsSqs:
		return c.publishSQS(ctx, publisher, message)
	case enum.PublisherProviderAwsSns:
		return c.publishSNS(ctx, publisher, message)
	default:
		return errors.Newf("Provider %s not supported", publisher.Provider())
	}
}

func (c client) publishSQS(ctx context.Context, publisher *vo.Publisher, message *vo.Message) error {
	if checker.IsNil(c.sqs) {
		return errors.New("SQS client not configuration. Please check your configuration.")
	}
	_, err := c.sqs.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:            converter.ToPointer(message.Body()),
		QueueUrl:               converter.ToPointer(publisher.Reference()),
		DelaySeconds:           int32(message.Delay().Time().Seconds()),
		MessageDeduplicationId: message.DeduplicationID(),
		MessageGroupId:         message.GroupID(),
	})
	return err
}

func (c client) publishSNS(ctx context.Context, publisher *vo.Publisher, message *vo.Message) error {
	if checker.IsNil(c.sns) {
		return errors.New("SNS client not configuration. Please check your configuration.")
	}
	_, err := c.sns.Publish(ctx, &sns.PublishInput{
		Message:                converter.ToPointer(message.Body()),
		MessageDeduplicationId: message.DeduplicationID(),
		MessageGroupId:         message.GroupID(),
		TopicArn:               converter.ToPointer(publisher.Reference()),
	})
	return err
}
