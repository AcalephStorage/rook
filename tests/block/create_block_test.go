/*
Copyright 2016 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package block

import (
	"strings"
	"testing"

	"github.com/rook/rook/pkg/model"
	"github.com/rook/rook/tests/framework/clients"
	"github.com/rook/rook/tests/framework/contracts"
	"github.com/rook/rook/tests/framework/enums"
	"github.com/rook/rook/tests/framework/installer"
	"github.com/rook/rook/tests/framework/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	defaultPool    = "rbd"
	pool1          = "rook_test_pool"
	blockImageName = "testImage"
)

func TestBlockCreate(t *testing.T) {
	suite.Run(t, new(BlockImageCreateSuite))
}

type BlockImageCreateSuite struct {
	suite.Suite
	testClient     *clients.TestClient
	rc             contracts.RestAPIOperator
	initBlockCount int
	installer      *installer.InstallHelper
}

func (s *BlockImageCreateSuite) SetupSuite() {

	kh, err := utils.CreatK8sHelper()
	require.NoError(s.T(), err)

	s.installer = installer.NewK8sRookhelper(kh.Clientset)

	err = s.installer.InstallRookOnK8s()
	require.NoError(s.T(), err)

	s.testClient, err = clients.CreateTestClient(enums.Kubernetes, kh)
	require.Nil(s.T(), err)

	s.rc = s.testClient.GetRestAPIClient()
	initialBlocks, err := s.rc.GetBlockImages()
	require.Nil(s.T(), err)
	s.initBlockCount = len(initialBlocks)
}

//Test Creating Block image on default pool(rbd)
func (s *BlockImageCreateSuite) TestCreatingNewBlockImageOnDefaultPool() {

	s.T().Log("Test Creating new block image  for default pool")
	newImage := model.BlockImage{Name: blockImageName, Size: 123, PoolName: defaultPool}
	cbi, err := s.rc.CreateBlockImage(newImage)
	require.Nil(s.T(), err)
	require.Contains(s.T(), cbi, "succeeded created image")
	b, _ := s.rc.GetBlockImages()
	require.Equal(s.T(), s.initBlockCount+1, len(b), "Make sure new block image is created")

}

//Test Creating Block image on custom pool
func (s *BlockImageCreateSuite) TestCreatingNewBlockImageOnCustomPool() {

	s.T().Log("Test Creating new block image for custom pool")
	newPool := model.Pool{Name: pool1}
	_, err := s.rc.CreatePool(newPool)
	require.Nil(s.T(), err)

	newImage := model.BlockImage{Name: blockImageName, Size: 123, PoolName: newPool.Name}
	cbi, err := s.rc.CreateBlockImage(newImage)
	require.Nil(s.T(), err)
	require.Contains(s.T(), cbi, "succeeded created image")
	b, _ := s.rc.GetBlockImages()
	require.Equal(s.T(), s.initBlockCount+1, len(b), "Make sure new block image is created")

}

//Test Creating Block image twice on same pool
func (s *BlockImageCreateSuite) TestRecreatingBlockImageForSamePool() {

	s.T().Log("Test Case when Block Image is created with Name that is already used by another block on same pool")
	// create new block image
	newImage := model.BlockImage{Name: blockImageName, Size: 123, PoolName: defaultPool}
	cbi, err := s.rc.CreateBlockImage(newImage)
	require.Nil(s.T(), err)
	require.Contains(s.T(), cbi, "succeeded created image")
	b, _ := s.rc.GetBlockImages()
	require.Equal(s.T(), s.initBlockCount+1, len(b), "Make sure new block image is created")

	//create same block again on same pool
	newImage2 := model.BlockImage{Name: blockImageName, Size: 2897, PoolName: defaultPool}
	cbi2, err := s.rc.CreateBlockImage(newImage2)
	require.Error(s.T(), err, "Make sure dupe block is not created")
	require.NotContains(s.T(), cbi2, "succeeded created image")
	b2, _ := s.rc.GetBlockImages()
	require.Equal(s.T(), len(b), len(b2), "Make sure new block image is not created")

}

//Test Creating Block image twice on different pool
func (s *BlockImageCreateSuite) TestRecreatingBlockImageForDifferentPool() {

	s.T().Log("Test Case when Block Image is created with Name that is already used by another block on different pool")
	// create new block image
	newImage := model.BlockImage{Name: blockImageName, Size: 123, PoolName: defaultPool}
	cbi, err := s.rc.CreateBlockImage(newImage)
	require.Nil(s.T(), err)
	require.Contains(s.T(), cbi, "succeeded created image")
	b, _ := s.rc.GetBlockImages()
	require.Equal(s.T(), s.initBlockCount+1, len(b), "Make sure new block image is created")

	newPool := model.Pool{Name: pool1}
	_, err = s.rc.CreatePool(newPool)
	require.Nil(s.T(), err)

	//create same block again on different pool
	newImage2 := model.BlockImage{Name: blockImageName, Size: 2897, PoolName: newPool.Name}
	cbi2, err := s.rc.CreateBlockImage(newImage2)
	require.Nil(s.T(), err)
	require.Contains(s.T(), cbi2, "succeeded created image")
	b2, _ := s.rc.GetBlockImages()
	require.Equal(s.T(), len(b)+1, len(b2), "Make sure new block image is created")

}

// Delete all Block images that have the word Test in their name
func (s *BlockImageCreateSuite) TearDownTest() {

	blocks, _ := s.rc.GetBlockImages()

	for _, b := range blocks {
		if strings.Contains(b.Name, "test") {
			s.rc.DeleteBlockImage(b)
		}
	}
}
func (s *BlockImageCreateSuite) TearDownSuite() {

	s.installer.UninstallRookFromK8s()

}
