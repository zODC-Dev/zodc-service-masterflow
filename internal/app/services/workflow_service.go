package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/filters"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type WorkflowService struct {
	db                 *sql.DB
	nodeRepo           *repositories.NodeRepository
	workflowRepo       *repositories.WorkflowRepository
	nodeConnectionRepo *repositories.NodeConnectionRepository
	nodeGroupRepo      *repositories.NodeGroupRepository
}

func NewWorkflowService(db *sql.DB, nodeRepo *repositories.NodeRepository, workflowRepo *repositories.WorkflowRepository, nodeConnectionRepository *repositories.NodeConnectionRepository) *WorkflowService {
	return &WorkflowService{
		db:                 db,
		nodeRepo:           nodeRepo,
		workflowRepo:       workflowRepo,
		nodeConnectionRepo: nodeConnectionRepository,
	}
}

func (s *WorkflowService) Create(ctx context.Context, req *requests.WorkflowRequest) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//Workflow Create
	workflowModel := model.Workflows{}
	if err := utils.Mapper(req, &workflowModel); err != nil {
		return err
	}

	workflow, err := s.workflowRepo.Create(ctx, tx, workflowModel)
	if err != nil {
		return err
	}

	//Node Create
	if len(req.Nodes) == 0 {
		return errors.New("requires at least 1 node")
	}

	nodesModel := []model.Nodes{}

	for i := range req.Nodes {
		node := req.Nodes[i]

		nodeModel := model.Nodes{
			WorkflowID: workflow.ID,
			X:          node.Position.X,
			Y:          node.Position.Y,
			Width:      node.Size.Width,
			Height:     node.Size.Height,
		}
		if err := utils.Mapper(node, &nodeModel); err != nil {
			return err
		}

		nodesModel = append(nodesModel, nodeModel)
	}

	err = s.nodeRepo.Create(ctx, tx, nodesModel)
	if err != nil {
		return err
	}

	//Connection Create
	if len(req.Connections) == 0 {
		return errors.New("requires at least 1 connection")
	}

	connectionsModel := []model.NodeConnections{}

	for i := range req.Connections {
		connectionReq := req.Connections[i]

		connectionModel := model.NodeConnections{
			ID:         connectionReq.Id,
			FromNodeID: connectionReq.From,
			ToNodeID:   connectionReq.To,
			Type:       connectionReq.Type,
			WorkflowID: workflow.ID,
		}

		connectionsModel = append(connectionsModel, connectionModel)
	}

	if err := s.nodeConnectionRepo.Create(ctx, tx, connectionsModel); err != nil {
		return err
	}

	//Group Create
	if len(req.Groups) != 0 {
		groupsModel := []model.NodeGroups{}

		for i := range req.Groups {
			groupReq := req.Groups[i]

			groupModel := model.NodeGroups{
				X:          groupReq.Position.X,
				Y:          groupReq.Position.Y,
				Width:      groupReq.Size.Width,
				Height:     groupReq.Size.Height,
				WorkflowID: workflow.ID,
			}

			if err := utils.Mapper(groupReq, &groupModel); err != nil {
				return err
			}

			groupsModel = append(groupsModel, groupModel)
		}

		if err := s.nodeGroupRepo.Create(ctx, tx, groupsModel); err != nil {
			return err
		}
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *WorkflowService) FindAll(ctx context.Context, workflowQueryParam *queryparams.WorkflowQueryParam) ([]responses.WorkflowResponse, error) {
	workflowsResponse := []responses.WorkflowResponse{}

	workflowFilter := filters.WorkflowFilter{}
	if err := utils.Mapper(&workflowQueryParam, &workflowFilter); err != nil {
		return workflowsResponse, err
	}

	workflows, err := s.workflowRepo.FindAll(ctx, s.db, workflowFilter)
	if err != nil {
		return workflowsResponse, err
	}

	for i := range workflows {
		workflowResponse := responses.WorkflowResponse{}
		connectionsResponse := []responses.ConnectionResponse{}
		nodesResponse := []responses.NodeResponse{}
		groupsResponse := []responses.GroupResponse{}

		for j := range workflows[i].Connections {
			connection := workflows[i].Connections[j]

			connectionResponse := responses.ConnectionResponse{
				Id:   connection.ID,
				From: connection.FromNodeID,
				To:   connection.ToNodeID,
				Type: connection.Type,
			}

			connectionsResponse = append(connectionsResponse, connectionResponse)
		}

		for j := range workflows[i].Nodes {
			node := workflows[i].Nodes[j]

			nodeResponse := responses.NodeResponse{
				Position: types.Position{
					X: node.X,
					Y: node.Y,
				},
				Size: types.Size{
					Width:  node.Width,
					Height: node.Height,
				},
			}
			if err := utils.Mapper(node, &nodeResponse); err != nil {
				return workflowsResponse, err
			}

			nodesResponse = append(nodesResponse, nodeResponse)
		}

		for j := range workflows[i].Groups {
			group := &workflows[i].Groups[j]

			if group.Type == nil {
				group.Type = new(string)
				*group.Type = ""
			}

			groupResponse := responses.GroupResponse{
				Position: types.Position{
					X: group.X,
					Y: group.Y,
				},
				Size: types.Size{
					Width:  group.Width,
					Height: group.Height,
				},
			}
			if err := utils.Mapper(group, &groupResponse); err != nil {
				return workflowsResponse, err
			}

			groupsResponse = append(groupsResponse, groupResponse)
		}

		if err := utils.Mapper(workflows[i], &workflowResponse); err != nil {
			return workflowsResponse, err
		}

		workflowResponse.Connections = connectionsResponse
		workflowResponse.Nodes = nodesResponse
		workflowResponse.Groups = groupsResponse

		workflowsResponse = append(workflowsResponse, workflowResponse)
	}

	return workflowsResponse, nil
}
