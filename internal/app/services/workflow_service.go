package services

import (
	"context"
	"database/sql"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type WorkflowService struct {
	db                       *sql.DB
	nodeRepo                 *repositories.NodeRepository
	workflowRepo             *repositories.WorkflowRepository
	nodeConnectionRepository *repositories.NodeConnectionRepository
}

func NewWorkflowService(db *sql.DB, nodeRepo *repositories.NodeRepository, workflowRepo *repositories.WorkflowRepository, nodeConnectionRepository *repositories.NodeConnectionRepository) *WorkflowService {
	return &WorkflowService{
		db:                       db,
		nodeRepo:                 nodeRepo,
		workflowRepo:             workflowRepo,
		nodeConnectionRepository: nodeConnectionRepository,
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
	nodesModel := []model.Nodes{}

	for i := range req.Nodes {
		node := req.Nodes[i]

		nodeModel := model.Nodes{}
		if err := utils.Mapper(node, &nodeModel); err != nil {
			return err
		}

		nodeModel.WorkflowID = workflow.ID

		nodesModel = append(nodesModel, nodeModel)
	}

	err = s.nodeRepo.Create(ctx, tx, nodesModel)
	if err != nil {
		return err
	}

	//Connection Create
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

	if err := s.nodeConnectionRepository.Create(ctx, tx, connectionsModel); err != nil {
		return err
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
