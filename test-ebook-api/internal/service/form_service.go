package service

import (
	"errors"
	"regexp"
	"test-ebook-api/internal/model"
	"test-ebook-api/internal/repository"
)

type FormService struct {
	repo         *repository.FormRepository
	categoryRepo *repository.StandardRepository
}

func NewFormService(repo *repository.FormRepository, catRepo *repository.StandardRepository) *FormService {
	return &FormService{repo: repo, categoryRepo: catRepo}
}

func (s *FormService) GetForms() ([]model.Form, error) {
	return s.repo.GetFormsWithFields()
}

func (s *FormService) GetFormByID(id uint) (*model.Form, error) {
	return s.repo.FindFormByID(id)
}

// GetOrCreateGlobalForm 获取或创建全局唯一属性表单
func (s *FormService) GetOrCreateGlobalForm() (*model.Form, error) {
	forms, err := s.repo.GetFormsWithFields()
	if err != nil {
		return nil, err
	}
	// 优先查找名为"全局属性"的 Form
	for i := range forms {
		if forms[i].Name == "全局属性" {
			return &forms[i], nil
		}
	}
	// 若只有一个 Form，直接将其作为全局属性
	if len(forms) == 1 {
		return &forms[0], nil
	}
	// 否则创建一个新的全局属性 Form
	form := &model.Form{
		Name:        "全局属性",
		Description: "系统全局文档属性配置",
	}
	if err := s.repo.CreateForm(form); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *FormService) CreateForm(name, description string) (*model.Form, error) {
	form := &model.Form{
		Name:        name,
		Description: description,
	}
	if err := s.repo.CreateForm(form); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *FormService) UpdateForm(id uint, name, description string) error {
	form, err := s.repo.FindFormByID(id)
	if err != nil {
		return errors.New("表单不存在")
	}
	form.Name = name
	form.Description = description
	return s.repo.UpdateForm(form)
}

func (s *FormService) DeleteForm(id uint) error {
	// 检查是否有分类在使用此表单
	count, err := s.categoryRepo.CountCategoriesByFormID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该表单模板正被分类使用中，无法删除")
	}
	return s.repo.DeleteForm(id)
}

func (s *FormService) SaveFormFields(formID uint, fields []model.FormField) error {
	if _, err := s.repo.FindFormByID(formID); err != nil {
		return errors.New("表单不存在")
	}

	// 校验字段标识合法性: 字母开头，仅由字母/数字/下划线组成
	keyRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	for _, f := range fields {
		if !keyRegex.MatchString(f.FieldKey) {
			return errors.New("字段标示 [" + f.FieldKey + "] 格式非法：必须以字母开头，且仅能包含字母、数字及下划线")
		}
	}

	return s.repo.UpdateFormFields(formID, fields)
}

func (s *FormService) BindCategoriesToForm(formID uint, categoryIDs []uint) error {
	if formID > 0 {
		if _, err := s.repo.FindFormByID(formID); err != nil {
			return errors.New("表单不存在")
		}
	}

	// 批量更新分类的 form_id
	// 直接使用 categoryRepo 的 DB 进行更新
	db := s.categoryRepo.GetDB()
	return db.Model(&model.Category{}).Where("id IN ?", categoryIDs).Update("form_id", formID).Error
}
