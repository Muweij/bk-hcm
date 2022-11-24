/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 混合云管理平台 (BlueKing - Hybrid Cloud Management System) available.
 * Copyright (C) 2022 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 *
 * to the current version of the project delivered to anyone in the future.
 */

package shutdown

import (
	"testing"
	"time"
)

func TestShutdownGracefully(t *testing.T) {
	job := func(name string) {
		notifier := AddNotifier()
		select {
		case <-notifier.Signal:
			t.Logf("%s received shutdown notify message.", name)
			notifier.Done()
		}
	}

	go job("tom")
	go job("jerry")
	time.Sleep(time.Second)

	RegisterFirstShutdown(func() {
		t.Logf("I'm the first shutdown job!!!")
	})

	SignalShutdownGracefully()

	WaitShutdown(0)
}

func TestShutdownTimeout(t *testing.T) {
	job := func(name string) {
		notifier := AddNotifier()
		select {
		case <-notifier.Signal:
			t.Logf("%s received shutdown notify message.", name)
			// sleep timeout.
			time.Sleep(6 * time.Second)
			notifier.Done()
		}
	}

	go job("tom")
	go job("jerry")
	time.Sleep(time.Second)

	t.Logf("send shutdown notify")
	SignalShutdownGracefully()

	time.Sleep(time.Second)
	t.Logf("wait job exec completed")
	WaitShutdown(20)
}
