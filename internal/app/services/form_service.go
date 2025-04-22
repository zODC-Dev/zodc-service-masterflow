package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type FormService struct {
	db         *sql.DB
	formRepo   *repositories.FormRepository
	natsClient *nats.NATSClient
}

func NewFormService(db *sql.DB, formRepo *repositories.FormRepository, natsClient *nats.NATSClient) *FormService {
	return &FormService{
		db:         db,
		formRepo:   formRepo,
		natsClient: natsClient,
	}
}

func (s *FormService) CreateFormTemplate(ctx context.Context, req *requests.FormTemplateCreate) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create Form Template
	formTemplate := model.FormTemplates{}
	if err := utils.Mapper(req, &formTemplate); err != nil {
		return fmt.Errorf("mapping form template failed: %w", err)
	}

	// Datasheet Mapping
	if req.DataSheet != nil {
		datasheet := string(*req.DataSheet)
		formTemplate.DataSheet = &datasheet
	}
	if req.TemplateID != nil {
		formTemplate.TemplateID = req.TemplateID
	}

	formTemplate.Type = "USER"
	formTemplate.Tag = "FORM"

	formTemplate, err = s.formRepo.CreateFormTemplate(ctx, tx, formTemplate)
	if err != nil {
		return fmt.Errorf("create form template failed: %w", err)
	}

	// Create Form Template Version
	formTemplateVersion := model.FormTemplateVersions{
		Version:        1,
		FormTemplateID: formTemplate.ID,
	}

	formTemplateVersion, err = s.formRepo.CreateFormTemplateVersion(ctx, tx, formTemplateVersion)
	if err != nil {
		return fmt.Errorf("create form template version failed: %w", err)
	}

	// Create Form Template Fields
	formTemplateFields := []model.FormTemplateFields{}
	for i := range req.FormFields {
		for j := range req.FormFields[i] {
			formTemplateField := model.FormTemplateFields{}

			if err := utils.Mapper(req.FormFields[i][j], &formTemplateField); err != nil {
				return fmt.Errorf("mapping form template field failed: %w", err)
			}

			// AdvancedOptions Mapping
			if req.FormFields[i][j].AdvancedOptions != nil {
				advancedOptions := string(*req.FormFields[i][j].AdvancedOptions)
				formTemplateField.AdvancedOptions = &advancedOptions
			}

			//
			formTemplateField.ColNum = int32(i)

			//
			formTemplateField.FormTemplateVersionID = formTemplateVersion.ID

			//
			formTemplateFields = append(formTemplateFields, formTemplateField)
		}

	}

	if err := s.formRepo.CreateFormTemplateFields(ctx, tx, formTemplateFields); err != nil {
		return fmt.Errorf("create form template field failed: %w", err)
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *FormService) FindAllFormTemplate(ctx context.Context, queryParam queryparams.FormQueryParam) ([]responses.FormTemplateFindAll, error) {
	formTemplatesResponse := []responses.FormTemplateFindAll{}

	formTemplates, err := s.formRepo.FindAllFormTemplate(ctx, s.db, queryParam)

	if err != nil {
		return formTemplatesResponse, err
	}

	for _, formTemplate := range formTemplates {

		formTemplateResponse := responses.FormTemplateFindAll{}
		if err := utils.Mapper(formTemplate, &formTemplateResponse); err != nil {
			return formTemplatesResponse, fmt.Errorf("map form template fail: %w", err)
		}

		if formTemplate.DataSheet != nil {
			var dataSheet map[string]interface{}
			if err := json.Unmarshal([]byte(*formTemplate.DataSheet), &dataSheet); err != nil {
				return formTemplatesResponse, fmt.Errorf("unmarshal data sheet fail: %w", err)
			}
			formTemplateResponse.DataSheet = &dataSheet
		}

		formTemplateResponse.Category = responses.CategoryFindAll{}
		if err := utils.Mapper(formTemplate.Category, &formTemplateResponse.Category); err != nil {
			return formTemplatesResponse, fmt.Errorf("map category fail: %w", err)
		}

		formTemplateResponse.Version = formTemplate.Version.Version

		formTemplatesResponse = append(formTemplatesResponse, formTemplateResponse)
	}

	return formTemplatesResponse, nil
}

func (s *FormService) FindAllFormTemplateFieldsByFormTemplateId(ctx context.Context, formTemplateId int32) ([][]responses.FormTemplateFieldsFindAll, error) {
	fieldsResponse := [][]responses.FormTemplateFieldsFindAll{}

	formTemplateFields, err := s.formRepo.FindAllFormTemplateFieldsByFormTemplateId(ctx, s.db, formTemplateId)
	if err != nil {
		return fieldsResponse, err
	}

	for _, formformTemplateField := range formTemplateFields {

		colIndex := formformTemplateField.ColNum

		for len(fieldsResponse) <= int(colIndex) {
			fieldsResponse = append(fieldsResponse, []responses.FormTemplateFieldsFindAll{})
		}

		fieldResponse := responses.FormTemplateFieldsFindAll{}

		//Mapping AdvancedOptions
		var advancedOptions map[string]interface{}
		if err := json.Unmarshal([]byte(*formformTemplateField.AdvancedOptions), &advancedOptions); err != nil {
			return fieldsResponse, err
		}
		fieldResponse.AdvancedOptions = advancedOptions
		if err := utils.Mapper(formformTemplateField, &fieldResponse); err != nil {
			return fieldsResponse, err
		}

		fieldsResponse[colIndex] = append(fieldsResponse[colIndex], fieldResponse)
	}

	return fieldsResponse, nil
}

func (s *FormService) FindOneFormTemplateDetailByFormTemplateId(ctx context.Context, formTemplateId int32) (responses.FormTemplateDetails, error) {
	formTemplateDetails := responses.FormTemplateDetails{}

	formTemplate, err := s.formRepo.FindOneFormTemplateByFormTemplateId(ctx, s.db, formTemplateId)
	if err != nil {
		return formTemplateDetails, err
	}

	formTemplateResponse := responses.FormTemplateFindAll{}
	if err := utils.Mapper(formTemplate, &formTemplateResponse); err != nil {
		return formTemplateDetails, fmt.Errorf("map form template fail: %w", err)
	}

	if formTemplate.DataSheet != nil {
		var dataSheet map[string]interface{}
		if err := json.Unmarshal([]byte(*formTemplate.DataSheet), &dataSheet); err != nil {
			return formTemplateDetails, fmt.Errorf("unmarshal data sheet fail: %w", err)
		}
		formTemplateResponse.DataSheet = &dataSheet
	}

	formTemplateResponse.Category = responses.CategoryFindAll{}
	if err := utils.Mapper(formTemplate.Category, &formTemplateResponse.Category); err != nil {
		return formTemplateDetails, fmt.Errorf("map category fail: %w", err)
	}

	formTemplateResponse.Version = formTemplate.Version.Version

	fieldsResponse := [][]responses.FormTemplateFieldsFindAll{}
	formTemplateFields, err := s.formRepo.FindAllFormTemplateFieldsByFormTemplateId(ctx, s.db, formTemplateId)
	if err != nil {
		return formTemplateDetails, err
	}

	for _, formformTemplateField := range formTemplateFields {

		colIndex := formformTemplateField.ColNum

		for len(fieldsResponse) <= int(colIndex) {
			fieldsResponse = append(fieldsResponse, []responses.FormTemplateFieldsFindAll{})
		}

		fieldResponse := responses.FormTemplateFieldsFindAll{}

		//Mapping AdvancedOptions
		var advancedOptions map[string]interface{}
		if err := json.Unmarshal([]byte(*formformTemplateField.AdvancedOptions), &advancedOptions); err != nil {
			return formTemplateDetails, err
		}
		fieldResponse.AdvancedOptions = advancedOptions
		if err := utils.Mapper(formformTemplateField, &fieldResponse); err != nil {
			return formTemplateDetails, err
		}

		fieldsResponse[colIndex] = append(fieldsResponse[colIndex], fieldResponse)
	}

	formTemplateDetails.Template = formTemplateResponse
	formTemplateDetails.Fields = fieldsResponse

	return formTemplateDetails, nil
}

func (s *FormService) UpdateFormTemplate(ctx context.Context, req *requests.FormTemplateUpdate, formTemplateId int32) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	formTemplate, err := s.formRepo.FindOneFormTemplateByFormTemplateId(ctx, s.db, formTemplateId)
	if err != nil {
		return err
	}

	if req.FileName != "" {
		formTemplate.FileName = req.FileName
	}

	if req.Title != "" {
		formTemplate.Title = req.Title
	}

	if req.CategoryID != nil {
		formTemplate.CategoryID = req.CategoryID
	}

	if req.TemplateID != nil {
		formTemplate.TemplateID = req.TemplateID
	}

	if req.DataSheet != nil {
		datasheet := string(*req.DataSheet)
		formTemplate.DataSheet = &datasheet
	}

	formTemplateModel := model.FormTemplates{}
	if err := utils.Mapper(formTemplate, &formTemplateModel); err != nil {
		return fmt.Errorf("mapping form template failed: %w", err)
	}

	if err := s.formRepo.UpdateFormTemplate(ctx, tx, formTemplateModel); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
