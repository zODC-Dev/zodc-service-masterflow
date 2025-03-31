package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type NodeService struct {
	DB              *sql.DB
	NodeRepo        *repositories.NodeRepository
	ConnectionRepo  *repositories.ConnectionRepository
	RequestRepo     *repositories.RequestRepository
	WorkflowService *WorkflowService
}

func NewNodeService(cfg NodeService) *NodeService {
	nodeService := &NodeService{}
	utils.Mapper(cfg, nodeService)
	return nodeService
}

// Function
func (s *NodeService) UpdateNodeStatusToInProcessing(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	node.Status = string(constants.NodeStatusInProgress)

	// Update Node
	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	// If Story Or Sub Workflow
	if err := s.WorkflowService.RunWorkflowIfItStoryOrSubWorkflow(ctx, tx, node); err != nil {
		return err
	}

	return nil
}

//

func (s *NodeService) StartNodeHandler(ctx context.Context, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//
	node, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return err
	}

	// Update Current Node Status To In Process
	node.Status = string(constants.NodeStatusInProgress)

	// Set actual start time
	now := time.Now()
	node.ActualStartTime = &now

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *NodeService) CompleteNodeHandler(ctx context.Context, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	node, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return err
	}

	if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
		return fmt.Errorf("story or sub workflow is auto complete by system, cant mark as complete by user")
	}

	// Update Current Node Status To Completed
	node.Status = string(constants.NodeStatusCompleted)

	// Set actual finish time
	now := time.Now()
	node.ActualEndTime = &now

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	connectionsToNode, err := s.ConnectionRepo.FindConnectionsWithToNodesByFromNodeId(ctx, s.DB, node.ID)
	if err != nil {
		return err
	}

	// Update Next Node To In Processing
	for i := range connectionsToNode {

		// Update connection to completeed
		connectionsToNode[i].IsCompleted = true
		connectionModel := model.Connections{}
		if err := utils.Mapper(connectionsToNode[i], &connectionModel); err != nil {
			return err
		}
		if err := s.ConnectionRepo.UpdateConnection(ctx, tx, connectionModel); err != nil {
			return err
		}

		// If Prevous Nodes not finish yet // If More than one node not completed then next node dont need to update status
		isUpdateNodeStatus := true

		connections, err := s.ConnectionRepo.FindConnectionsByToNodeId(ctx, s.DB, connectionsToNode[i].Node.ID)
		if err != nil {
			return err
		}
		for j := range connections {
			if !connections[j].IsCompleted {
				isUpdateNodeStatus = false
				break
			}
		}

		if isUpdateNodeStatus {
			// If Node is End Node
			if connectionsToNode[i].Node.Type == string(constants.NodeTypeEnd) {
				// Update end node to completed
				connectionsToNode[i].Node.Status = string(constants.NodeStatusCompleted)
				if err := s.NodeRepo.UpdateNode(ctx, tx, connectionsToNode[i].Node); err != nil {
					return err
				}

				// Mark request completed
				request, err := s.RequestRepo.FindRequestByNodeId(ctx, s.DB, connectionsToNode[i].Node.ID)
				if err != nil {
					return err
				}
				request.Status = string(constants.RequestStatusCompleted)
				if err := s.RequestRepo.UpdateRequest(ctx, tx, request); err != nil {
					return err
				}
			} else {
				if err := s.UpdateNodeStatusToInProcessing(ctx, tx, connectionsToNode[i].Node); err != nil {
					return err
				}
			}
		}
	}

	// Calculate Request Process
	request, _ := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, node.RequestID)
	totalCompletedNode := 0
	totalNode := len(request.Nodes)
	for _, requestNode := range request.Nodes {
		if requestNode.Status == string(constants.NodeStatusCompleted) {
			totalCompletedNode++
		}

		if requestNode.Type == string(constants.NodeTypeStart) || requestNode.Type == string(constants.NodeTypeEnd) {
			totalNode--
		}
	}
	if totalNode == 0 {
		request.Progress = 100
	} else {
		request.Progress = float32(float64(totalCompletedNode) / float64(totalNode) * 100)
	}
	requestModel := model.Requests{}
	if err := utils.Mapper(request, &requestModel); err != nil {
		return err
	}
	if err := s.RequestRepo.UpdateRequest(ctx, tx, requestModel); err != nil {
		return err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}
