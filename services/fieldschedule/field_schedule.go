package services

import (
	"context"
	"field-services/common/utils"
	"field-services/constants"
	constErrFieldSchedule "field-services/constants/errors/fieldschedule"
	"field-services/domains/dto"
	"field-services/domains/models"
	"field-services/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FieldScheduleService struct {
	repository repositories.IRepositoryRegistry
}

// Create implements [IFieldScheduleService].
func (f *FieldScheduleService) Create(ctx context.Context, req *dto.FieldScheduleRequest) error {
	field, err := f.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	fieldSchedules := make([]models.FieldSchedule, 0, len(req.FieldID))
	dateParsed, _ := time.Parse(time.DateOnly, req.Date)
	for _, timeID := range req.TimeIDs {
		scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, timeID)
		if err != nil {
			return err
		}

		schedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
			ctx,
			req.Date,
			int(scheduleTime.ID),
			int(field.ID),
		)
		if err != nil {
			return err
		}
		if schedule != nil {
			return constErrFieldSchedule.ErrFieldScheduleIsExist
		}

		fieldSchedules = append(fieldSchedules, models.FieldSchedule{
			UUID:    uuid.New(),
			FieldID: field.ID,
			TimeID:  scheduleTime.ID,
			Date:    dateParsed,
			Status:  constants.Available,
		})

	}

	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err
	}

	return nil
}

// Delete implements [IFieldScheduleService].
func (f *FieldScheduleService) Delete(ctx context.Context, uuid string) error {
	_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = f.repository.GetFieldSchedule().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

// GenerateScheduleForOneMonth implements [IFieldScheduleService].
func (f *FieldScheduleService) GenerateScheduleForOneMonth(ctx context.Context, req *dto.GenerateFieldScheduleForOneMonthRequest) error {
	field, err := f.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	times, err := f.repository.GetTime().FindAll(ctx)
	if err != nil {
		return err
	}

	numberOfDays := 30
	fieldSchedules := make([]models.FieldSchedule, 0, numberOfDays)
	now := time.Now().Add(time.Duration(1) * 24 * time.Hour)
	for i := 0; i < numberOfDays; i++ {
		currentDate := now.AddDate(0, 0, i)
		for _, item := range times {
			schedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
				ctx,
				currentDate.Format(time.DateOnly),
				int(item.ID),
				int(field.ID),
			)
			if err != nil {
				return err
			}

			if schedule != nil {
				return constErrFieldSchedule.ErrFieldScheduleIsExist
			}

			fieldSchedules = append(fieldSchedules, models.FieldSchedule{
				UUID:    uuid.New(),
				FieldID: field.ID,
				TimeID:  item.ID,
				Date:    currentDate,
				Status:  constants.Available,
			})
		}
	}

	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err
	}

	return nil
}

// GetAllByFieldIDAndDate implements [IFieldScheduleService].
func (f *FieldScheduleService) GetAllByFieldIDAndDate(
	ctx context.Context,
	uuid string,
	date string,
) ([]dto.FieldScheduleForBookingResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedules, err := f.repository.GetFieldSchedule().FindAllByFieldIDAndDate(ctx, int(field.ID), date)
	if err != nil {
		return nil, err
	}

	fieldSchedulesResults := make([]dto.FieldScheduleForBookingResponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		pricePerHour := float64(schedule.Field.PricePerHour)
		fieldSchedulesResults = append(fieldSchedulesResults, dto.FieldScheduleForBookingResponse{
			UUID:         schedule.UUID,
			Date:         f.convertMonthName(schedule.Date.Format(time.DateOnly)),
			Time:         schedule.Time.StartTime,
			Status:       schedule.Status.GetStatusString(),
			PricePerHour: utils.RupiahFormarter(&pricePerHour),
		})
	}
	return fieldSchedulesResults, nil
}

// GetAllWithPagination implements [IFieldScheduleService].
func (f *FieldScheduleService) GetAllWithPagination(
	ctx context.Context,
	param *dto.FieldScheduleRequestParam,
) (*utils.PaginationResult, error) {
	fieldSchedules, total, err := f.repository.GetFieldSchedule().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldSchedulesResults := make([]dto.FieldScheduleResponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		fieldSchedulesResults = append(fieldSchedulesResults, dto.FieldScheduleResponse{
			UUID:         schedule.UUID,
			FieldName:    schedule.Field.Name,
			Date:         schedule.Date.Format(time.DateOnly),
			PricePerHour: schedule.Field.PricePerHour,
			Status:       schedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", schedule.Time.StartTime, schedule.Time.EndTime),
			CreatedAt:    schedule.CreatedAt,
			UpdatedAt:    schedule.UpdatedAt,
		})
	}

	pagination := &utils.PaginationParams{
		Count: int(total),
		Limit: param.Limit,
		Page:  param.Page,
		Data:  fieldSchedulesResults,
	}

	response := utils.GeneratePagination(*pagination)
	return &response, nil
}

// GetByUUID implements [IFieldScheduleService].
func (f *FieldScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error) {
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleResponse{
		UUID:         fieldSchedule.UUID,
		FieldName:    fieldSchedule.Field.Name,
		PricePerHour: fieldSchedule.Field.PricePerHour,
		Date:         fieldSchedule.Date.Format(time.DateOnly),
		Status:       fieldSchedule.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", fieldSchedule.Time.StartTime, fieldSchedule.Time.EndTime),
		CreatedAt:    fieldSchedule.CreatedAt,
		UpdatedAt:    fieldSchedule.UpdatedAt,
	}

	return &response, nil
}

// Update implements [IFieldScheduleService].
func (f *FieldScheduleService) Update(ctx context.Context, uuid string, req *dto.UpdateFieladScheduleRequest) (*dto.FieldScheduleResponse, error) {
	fieldSchedules, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, req.TimeID)
	if err != nil {
		return nil, err
	}

	isTimeExist, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
		ctx,
		req.Date,
		int(scheduleTime.ID),
		int(fieldSchedules.FieldID),
	)
	if err != nil {
		return nil, err
	}

	if isTimeExist != nil && req.Date != fieldSchedules.Date.Format(time.DateOnly) {
		checkDate, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(
			ctx,
			req.Date,
			int(scheduleTime.ID),
			int(fieldSchedules.FieldID),
		)

		if err != nil {
			return nil, err
		}

		if checkDate != nil {
			return nil, constErrFieldSchedule.ErrFieldScheduleIsExist
		}
	}

	dateParsed, _ := time.Parse(time.DateOnly, req.Date)
	fieldResults, err := f.repository.GetFieldSchedule().Update(ctx, uuid, &models.FieldSchedule{
		Date:   dateParsed,
		TimeID: scheduleTime.ID,
	})

	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleResponse{
		UUID:         fieldSchedules.UUID,
		FieldName:    fieldResults.Field.Name,
		Date:         fieldResults.Date.Format(time.DateOnly),
		PricePerHour: fieldResults.Field.PricePerHour,
		Status:       fieldResults.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", scheduleTime.StartTime, scheduleTime.EndTime),
	}

	return &response, nil
}

// UpdateStatus implements [IFieldScheduleService].
func (f *FieldScheduleService) UpdateStatus(ctx context.Context, req *dto.UpdateStatusFieladScheduleRequest) error {
	for _, item := range req.FieldScheduleIDs {
		_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, item)
		if err != nil {
			return err
		}

		err = f.repository.GetFieldSchedule().UpdateStatus(ctx, item, constants.Booked)
		if err != nil {
			return err
		}
	}

	return nil
}

type IFieldScheduleService interface {
	GetAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) (*utils.PaginationResult, error)
	GetAllByFieldIDAndDate(context.Context, string, string) ([]dto.FieldScheduleForBookingResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldScheduleResponse, error)
	GenerateScheduleForOneMonth(context.Context, *dto.GenerateFieldScheduleForOneMonthRequest) error
	Create(context.Context, *dto.FieldScheduleRequest) error
	Update(context.Context, string, *dto.UpdateFieladScheduleRequest) (*dto.FieldScheduleResponse, error)
	UpdateStatus(context.Context, *dto.UpdateStatusFieladScheduleRequest) error
	Delete(context.Context, string) error
}

func NewFieldScheduleRegistry(repository repositories.IRepositoryRegistry) IFieldScheduleService {
	return &FieldScheduleService{
		repository: repository,
	}
}

func (f *FieldScheduleService) convertMonthName(inputDate string) string {
	date, err := time.Parse(time.DateOnly, inputDate)
	if err != nil {
		return ""
	}

	indonesianMonth := map[string]string{
		"Jan": "Jan",
		"Feb": "Feb",
		"Mar": "Mar",
		"Apr": "Apr",
		"May": "Mei",
		"Jun": "Jun",
		"Jul": "Jul",
		"Aug": "Aug",
		"Sep": "Sep",
		"Oct": "Oct",
		"Nov": "Nov",
		"Dec": "Dec",
	}
	formattedDate := date.Format("02 Jan")
	day := formattedDate[:3]
	month := formattedDate[3:]
	formattedDate = fmt.Sprintf("%s %s", day, indonesianMonth[month])
	return formattedDate
}
