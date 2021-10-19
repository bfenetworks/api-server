// Copyright (c) 2021 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ibasic

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type ExtraFile struct {
	ID          int64
	Name        string
	ProductID   int64
	Description string
	Md5         []byte
	Content     []byte
}

type ExtraFileParam struct {
	Name        *string
	Description *string
	Md5         []byte
	Content     []byte
}

type ExtraFileFilter struct {
	Name  *string
	Names []string
}

type ExtraFileStorager interface {
	CreateExtraFile(context.Context, *Product, ...*ExtraFileParam) error
	DeleteExtraFile(context.Context, *ExtraFileFilter) error
	FetchExtraFiles(context.Context, *ExtraFileFilter) ([]*ExtraFile, error)
}

func ExtraFilePath(moduleDir string, product *Product, fileName string) string {
	fileName = strings.ReplaceAll(fileName, string(os.PathSeparator), "_")
	return fmt.Sprintf("%s/%s/%s", moduleDir, strings.ToLower(product.Name), fileName)
}

type ExtraFileManager struct {
	storager ExtraFileStorager
}

func NewExtraFileManager(storager ExtraFileStorager) *ExtraFileManager {
	return &ExtraFileManager{
		storager: storager,
	}
}

func (em *ExtraFileManager) FetchExtraFile(ctx context.Context, fileName string) (*ExtraFile, error) {
	list, err := em.storager.FetchExtraFiles(ctx, &ExtraFileFilter{
		Name: &fileName,
	})
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}
