package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
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
		formTemplateResponse.VersionId = formTemplate.Version.ID

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
	for _, formformTemplateField := range formTemplate.Fields {

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

func (s *FormService) FindOneFormTemplateDetailByFormTemplateVersionId(ctx context.Context, formTemplateVersionId int32) (responses.FormTemplateDetails, error) {
	formTemplateDetails := responses.FormTemplateDetails{}

	formTemplate, err := s.formRepo.FindOneFormTemplateByFormTemplateVersionId(ctx, s.db, formTemplateVersionId)
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
	for _, formformTemplateField := range formTemplate.Fields {

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

	formTemplateModel := model.FormTemplates{
		ID:             formTemplate.ID,
		FileName:       req.FileName,
		Title:          req.Title,
		CurrentVersion: formTemplate.CurrentVersion,
		CategoryID:     req.CategoryID,
		TemplateID:     req.TemplateID,
		Description:    req.Description,
		Decoration:     req.Decoration,
		Tag:            formTemplate.Tag,
		Type:           formTemplate.Type,
	}

	if req.DataSheet != nil {
		datasheet := string(*req.DataSheet)
		formTemplateModel.DataSheet = &datasheet
	} else {
		formTemplateModel.DataSheet = formTemplate.DataSheet
	}

	if err := s.formRepo.UpdateFormTemplate(ctx, tx, formTemplateModel); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *FormService) ConfigFormTemplate(ctx context.Context, formTemplateId int32, req *[][]requests.FormTemplateFieldsCreate) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	formTemplate, err := s.formRepo.FindOneFormTemplateByFormTemplateId(ctx, s.db, formTemplateId)
	if err != nil {
		return err
	}

	formTemplate.CurrentVersion = formTemplate.CurrentVersion + 1
	formTemplateModel := model.FormTemplates{}
	if err := utils.Mapper(formTemplate, &formTemplateModel); err != nil {
		return fmt.Errorf("mapping form template failed: %w", err)
	}
	if err := s.formRepo.UpdateFormTemplate(ctx, tx, formTemplateModel); err != nil {
		return err
	}

	formTemplateVersion := model.FormTemplateVersions{
		Version:        formTemplate.CurrentVersion,
		FormTemplateID: formTemplate.ID,
	}
	formTemplateVersion, err = s.formRepo.CreateFormTemplateVersion(ctx, tx, formTemplateVersion)
	if err != nil {
		return fmt.Errorf("create form template version failed: %w", err)
	}

	formTemplateFields := []model.FormTemplateFields{}
	for i := range *req {
		for j := range (*req)[i] {
			formTemplateField := model.FormTemplateFields{}

			if err := utils.Mapper((*req)[i][j], &formTemplateField); err != nil {
				return fmt.Errorf("mapping form template field failed: %w", err)
			}

			// AdvancedOptions Mapping
			if (*req)[i][j].AdvancedOptions != nil {
				advancedOptions := string(*(*req)[i][j].AdvancedOptions)
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

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *FormService) FindOneFormDataByFormDataId(ctx context.Context, formDataId string) (responses.JiraFormDetailResponse, error) {
	formDataResponse := responses.JiraFormDetailResponse{}

	formData, err := s.formRepo.FindFormDataById(ctx, s.db, formDataId)
	if err != nil {
		return formDataResponse, err
	}

	formTemplate, err := s.FindOneFormTemplateDetailByFormTemplateId(ctx, formData.FormTemplates.ID)
	if err != nil {
		return responses.JiraFormDetailResponse{}, err
	}

	fieldMap := map[int32]string{}
	for _, formTemplateField := range formData.FormTemplateFields {
		fieldMap[formTemplateField.ID] = formTemplateField.FieldID
	}

	formDatas := []responses.NodeFormDataResponse{}
	for _, formData := range formData.FormFieldData {
		formDatas = append(formDatas, responses.NodeFormDataResponse{
			FieldId: fieldMap[formData.FormTemplateFieldID],
			Value:   formData.Value,
		})
	}

	formDataResponse.Template = formTemplate.Template
	formDataResponse.Fields = formTemplate.Fields
	formDataResponse.Data = formDatas

	return formDataResponse, nil
}

func (s *FormService) GetEditProfileFormTemplate(ctx context.Context) (responses.FormTemplateDetails, error) {
	formTemplate, err := s.FindOneFormTemplateDetailByFormTemplateId(ctx, constants.FormTemplateIDEditProfileForm)
	if err != nil {
		return responses.FormTemplateDetails{}, err
	}

	return formTemplate, nil
}

func (s *FormService) GetPerformanceEvaluateFormTemplate(ctx context.Context) (responses.FormTemplateDetails, error) {
	formTemplate, err := s.FindOneFormTemplateDetailByFormTemplateId(ctx, constants.FormTemplateIDPerformanceEvaluate)
	if err != nil {
		return responses.FormTemplateDetails{}, err
	}

	return formTemplate, nil
}
