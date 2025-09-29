package resolver

import (
	"fmt"
	"time"

	"github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	equipment            []*models.Equipment
	facilities           []*models.Facility
	maintenanceRecords   []*models.MaintenanceRecord
	operationalMetrics   []*models.OperationalMetric
	alerts               []*models.EquipmentAlert
	personnel            []models.Personnel
	managers             []*models.Manager
	associates           []*models.Associate
	maintenanceReminders []*models.MaintenanceReminder
	statusUpdates        []*models.StatusUpdate
	performanceAlerts    []*models.PerformanceAlert
}

// NewResolver creates a new resolver with sample data
func NewResolver() *Resolver {
	r := &Resolver{
		equipment:            make([]*models.Equipment, 0),
		facilities:           make([]*models.Facility, 0),
		maintenanceRecords:   make([]*models.MaintenanceRecord, 0),
		operationalMetrics:   make([]*models.OperationalMetric, 0),
		alerts:               make([]*models.EquipmentAlert, 0),
		maintenanceReminders: make([]*models.MaintenanceReminder, 0),
		statusUpdates:        make([]*models.StatusUpdate, 0),
		performanceAlerts:    make([]*models.PerformanceAlert, 0),
	}

	// Initialize with sample data
	r.initializeSampleData()
	return r
}

func (r *Resolver) initializeSampleData() {
	// Create sample facilities
	facility1 := &models.Facility{
		ID:          "facility-1",
		Name:        "Main Production Plant",
		Address:     "123 Industrial Blvd, Detroit, MI 48201, USA",
		Location:    &models.Location{Latitude: 42.3314, Longitude: -83.0458, Altitude: &[]float64{200.0}[0]},
		Status:      models.FacilityOperationalStatusOperational,
		ContactInfo: &models.ContactInfo{Phone: "+1-313-555-0100", Email: "plant@company.com"},
		Capacity:    1000,
		Utilization: 75.5,
	}

	facility2 := &models.Facility{
		ID:          "facility-2",
		Name:        "Secondary Assembly Line",
		Address:     "456 Manufacturing St, Chicago, IL 60601, USA",
		Location:    &models.Location{Latitude: 41.8781, Longitude: -87.6298, Altitude: &[]float64{180.0}[0]},
		Status:      models.FacilityOperationalStatusOperational,
		ContactInfo: &models.ContactInfo{Phone: "+1-312-555-0200", Email: "assembly@company.com"},
		Capacity:    500,
		Utilization: 60.0,
	}

	r.facilities = append(r.facilities, facility1, facility2)

	// Create sample equipment
	equipment1 := &models.Equipment{
		ID:           "equipment-1",
		Name:         "CNC Milling Machine Alpha",
		Description:  "High-precision CNC milling machine for precision parts manufacturing",
		Manufacturer: "Haas Automation",
		Model:        "VF-2SS",
		SerialNumber: "HAAS-2020-001",
		Type:         models.EquipmentTypeCncMill,
		Status:       models.EquipmentStatusRunning,
		Facility:     facility1,
		Specifications: &models.EquipmentSpecifications{
			Dimensions:           &models.Dimensions{Length: 2000, Width: 1500, Height: 1800},
			Weight:               2500,
			PowerConsumption:     15.5,
			MaxSpeed:             5000,
			OperatingTemperature: &models.TemperatureRange{Min: -10, Max: 50},
			ElectricalSpecs: &models.ElectricalSpecs{
				Voltage:   480,
				Frequency: 60,
			},
			EnvironmentalRequirements: []string{"Clean room environment", "Temperature controlled"},
		},
		InstalledAt:         "2020-03-15",
		LastMaintenanceAt:   &[]string{"2024-09-01"}[0],
		NextMaintenanceAt:   &[]string{"2024-12-01"}[0],
		Efficiency:          92.5,
		TotalOperatingHours: 8760,
	}

	equipment2 := &models.Equipment{
		ID:           "equipment-2",
		Name:         "Industrial Conveyor System Beta",
		Description:  "Main conveyor system for assembly line operations",
		Manufacturer: "Flexco",
		Model:        "FC-5000",
		SerialNumber: "FLEX-2021-002",
		Type:         models.EquipmentTypeConveyorBelt,
		Status:       models.EquipmentStatusMaintenance,
		Facility:     facility2,
		Specifications: &models.EquipmentSpecifications{
			Dimensions:           &models.Dimensions{Length: 5000, Width: 800, Height: 1200},
			Weight:               1200,
			PowerConsumption:     7.5,
			MaxSpeed:             2.5,
			OperatingTemperature: &models.TemperatureRange{Min: 0, Max: 40},
			ElectricalSpecs: &models.ElectricalSpecs{
				Voltage:   240,
				Frequency: 60,
			},
			EnvironmentalRequirements: []string{"Standard industrial environment"},
		},
		InstalledAt:         "2021-08-20",
		LastMaintenanceAt:   &[]string{"2024-08-15"}[0],
		NextMaintenanceAt:   &[]string{"2024-10-15"}[0],
		Efficiency:          78.0,
		TotalOperatingHours: 5256,
	}

	r.equipment = append(r.equipment, equipment1, equipment2)

	// Create sample maintenance records
	maintenance1 := &models.MaintenanceRecord{
		ID:                 "maintenance-1",
		Equipment:          equipment1,
		Type:               models.MaintenanceTypePreventive,
		Status:             models.MaintenanceStatusCompleted,
		Priority:           models.MaintenancePriorityMedium,
		ScheduledDate:      "2024-09-01",
		CompletedDate:      &[]string{"2024-09-01"}[0],
		Description:        "Routine preventive maintenance - lubrication and calibration",
		AssignedTechnician: "John Smith",
		EstimatedDuration:  4,
	}

	maintenance2 := &models.MaintenanceRecord{
		ID:                 "maintenance-2",
		Equipment:          equipment2,
		Type:               models.MaintenanceTypeCorrective,
		Status:             models.MaintenanceStatusScheduled,
		Priority:           models.MaintenancePriorityHigh,
		ScheduledDate:      "2024-10-15",
		Description:        "Belt replacement and motor inspection",
		AssignedTechnician: "Jane Doe",
		EstimatedDuration:  8,
	}

	r.maintenanceRecords = append(r.maintenanceRecords, maintenance1, maintenance2)

	// Create sample operational metrics
	metric1 := &models.OperationalMetric{
		ID:          "metric-1",
		Equipment:   equipment1,
		MetricType:  models.MetricTypeTemperature,
		Value:       25.5,
		Unit:        "°C",
		RecordedAt:  time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		TargetValue: &[]float64{30.0}[0],
	}

	metric2 := &models.OperationalMetric{
		ID:          "metric-2",
		Equipment:   equipment1,
		MetricType:  models.MetricTypeVibration,
		Value:       0.8,
		Unit:        "mm/s",
		RecordedAt:  time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
		TargetValue: &[]float64{1.0}[0],
	}

	r.operationalMetrics = append(r.operationalMetrics, metric1, metric2)

	// Create sample alerts
	alert1 := &models.EquipmentAlert{
		ID:           "alert-1",
		Equipment:    equipment2,
		Type:         models.AlertTypeMaintenanceDue,
		Severity:     models.AlertSeverityMedium,
		Description:  "Belt tension below recommended threshold",
		GeneratedAt:  time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		Acknowledged: false,
	}

	r.alerts = append(r.alerts, alert1)

	// Create sample personnel
	manager1 := &models.Manager{
		ID:            "manager-1",
		Name:          "John Smith",
		Email:         "john.smith@company.com",
		Phone:         "+1-313-555-0101",
		JoinedAt:      "2020-01-15",
		Status:        models.PersonnelStatusActive,
		Department:    "Production",
		DirectReports: 5,
		Level:         4,
	}

	manager2 := &models.Manager{
		ID:            "manager-2",
		Name:          "Sarah Johnson",
		Email:         "sarah.johnson@company.com",
		Phone:         "+1-312-555-0201",
		JoinedAt:      "2019-06-01",
		Status:        models.PersonnelStatusActive,
		Department:    "Quality Control",
		DirectReports: 3,
		Level:         3,
	}

	associate1 := &models.Associate{
		ID:         "associate-1",
		Name:       "Mike Wilson",
		Email:      "mike.wilson@company.com",
		Phone:      "+1-313-555-0102",
		JoinedAt:   "2021-03-10",
		Status:     models.PersonnelStatusActive,
		JobTitle:   "Machine Operator",
		Department: "Production",
		ReportsTo:  manager1,
	}

	associate2 := &models.Associate{
		ID:         "associate-2",
		Name:       "Lisa Brown",
		Email:      "lisa.brown@company.com",
		Phone:      "+1-312-555-0202",
		JoinedAt:   "2022-01-20",
		Status:     models.PersonnelStatusActive,
		JobTitle:   "Quality Inspector",
		Department: "Quality Control",
		ReportsTo:  manager2,
	}

	associate3 := &models.Associate{
		ID:         "associate-3",
		Name:       "David Lee",
		Email:      "david.lee@company.com",
		Phone:      "+1-313-555-0103",
		JoinedAt:   "2020-11-05",
		Status:     models.PersonnelStatusActive,
		JobTitle:   "Maintenance Technician",
		Department: "Maintenance",
		ReportsTo:  manager1,
	}

	// Add personnel to global list
	r.personnel = append(r.personnel, manager1, manager2, associate1, associate2, associate3)

	// Add personnel to facilities
	facility1.Personnel = append(facility1.Personnel, manager1, associate1, associate3)
	facility2.Personnel = append(facility2.Personnel, manager2, associate2)

	// Create sample maintenance reminders
	reminder1 := &models.MaintenanceReminder{
		ID:            "reminder-1",
		Equipment:     equipment1,
		Type:          models.MaintenanceReminderTypeUpcomingMaintenance,
		Priority:      models.MaintenancePriorityMedium,
		Description:   "Scheduled preventive maintenance due in 2 weeks",
		CreatedAt:     time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		ScheduledDate: "2024-12-01",
		Acknowledged:  false,
	}

	reminder2 := &models.MaintenanceReminder{
		ID:             "reminder-2",
		Equipment:      equipment2,
		Type:           models.MaintenanceReminderTypeOverdueMaintenance,
		Priority:       models.MaintenancePriorityHigh,
		Description:    "Belt replacement maintenance is overdue",
		CreatedAt:      time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		ScheduledDate:  "2024-10-15",
		Acknowledged:   true,
		AcknowledgedAt: &[]string{time.Now().Add(-1 * time.Hour).Format(time.RFC3339)}[0],
		AcknowledgedBy: &[]string{"Jane Doe"}[0],
	}

	r.maintenanceReminders = append(r.maintenanceReminders, reminder1, reminder2)

	// Create sample status updates
	statusUpdate1 := &models.StatusUpdate{
		ID:             "status-1",
		Equipment:      equipment1,
		PreviousStatus: models.EquipmentStatusStopped,
		NewStatus:      models.EquipmentStatusRunning,
		Description:    "Equipment started after routine maintenance",
		ChangedAt:      time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
		ChangedBy:      &[]string{"Mike Wilson"}[0],
		Notes:          &[]string{"All systems operational, ready for production"}[0],
	}

	statusUpdate2 := &models.StatusUpdate{
		ID:             "status-2",
		Equipment:      equipment2,
		PreviousStatus: models.EquipmentStatusRunning,
		NewStatus:      models.EquipmentStatusMaintenance,
		Description:    "Equipment taken offline for scheduled maintenance",
		ChangedAt:      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		ChangedBy:      &[]string{"David Lee"}[0],
		Notes:          &[]string{"Belt replacement and motor inspection in progress"}[0],
	}

	r.statusUpdates = append(r.statusUpdates, statusUpdate1, statusUpdate2)

	// Create sample performance alerts
	perfAlert1 := &models.PerformanceAlert{
		ID:            "perf-alert-1",
		Equipment:     equipment1,
		MetricType:    models.MetricTypeEfficiency,
		CurrentValue:  85.2,
		ExpectedValue: 90.0,
		Threshold:     88.0,
		Severity:      models.AlertSeverityMedium,
		Description:   "Equipment efficiency below expected threshold",
		GeneratedAt:   time.Now().Add(-45 * time.Minute).Format(time.RFC3339),
		Acknowledged:  false,
	}

	perfAlert2 := &models.PerformanceAlert{
		ID:             "perf-alert-2",
		Equipment:      equipment2,
		MetricType:     models.MetricTypeVibration,
		CurrentValue:   2.5,
		ExpectedValue:  1.0,
		Threshold:      2.0,
		Severity:       models.AlertSeverityHigh,
		Description:    "Excessive vibration detected in conveyor system",
		GeneratedAt:    time.Now().Add(-15 * time.Minute).Format(time.RFC3339),
		Acknowledged:   true,
		AcknowledgedAt: &[]string{time.Now().Add(-5 * time.Minute).Format(time.RFC3339)}[0],
		AcknowledgedBy: &[]string{"Sarah Johnson"}[0],
	}

	r.performanceAlerts = append(r.performanceAlerts, perfAlert1, perfAlert2)
}

// Helper methods for finding entities by ID

func (r *Resolver) findEquipmentByID(id string) *models.Equipment {
	for _, equipment := range r.equipment {
		if equipment.ID == id {
			return equipment
		}
	}
	return nil
}

func (r *Resolver) findFacilityByID(id string) *models.Facility {
	for _, facility := range r.facilities {
		if facility.ID == id {
			return facility
		}
	}
	return nil
}

func (r *Resolver) findMaintenanceRecordByID(id string) *models.MaintenanceRecord {
	for _, maintenance := range r.maintenanceRecords {
		if maintenance.ID == id {
			return maintenance
		}
	}
	return nil
}

func (r *Resolver) findOperationalMetricByID(id string) *models.OperationalMetric {
	for _, metric := range r.operationalMetrics {
		if metric.ID == id {
			return metric
		}
	}
	return nil
}

func (r *Resolver) findAlertByID(id string) *models.EquipmentAlert {
	for _, alert := range r.alerts {
		if alert.ID == id {
			return alert
		}
	}
	return nil
}

// Helper methods for removing entities

func (r *Resolver) removeEquipmentFromFacility(facility *models.Facility, equipment *models.Equipment) {
	for i, e := range facility.Equipment {
		if e.ID == equipment.ID {
			facility.Equipment = append(facility.Equipment[:i], facility.Equipment[i+1:]...)
			break
		}
	}
}

func (r *Resolver) removeMaintenanceRecordsForEquipment(equipmentID string) {
	var filtered []*models.MaintenanceRecord
	for _, maintenance := range r.maintenanceRecords {
		if maintenance.Equipment.ID != equipmentID {
			filtered = append(filtered, maintenance)
		}
	}
	r.maintenanceRecords = filtered
}

func (r *Resolver) removeOperationalMetricsForEquipment(equipmentID string) {
	var filtered []*models.OperationalMetric
	for _, metric := range r.operationalMetrics {
		if metric.Equipment.ID != equipmentID {
			filtered = append(filtered, metric)
		}
	}
	r.operationalMetrics = filtered
}

func (r *Resolver) removeAlertsForEquipment(equipmentID string) {
	var filtered []*models.EquipmentAlert
	for _, alert := range r.alerts {
		if alert.Equipment.ID != equipmentID {
			filtered = append(filtered, alert)
		}
	}
	r.alerts = filtered
}

// GetEquipmentNotifications returns all equipment notifications (union type)
func (r *Resolver) GetEquipmentNotifications() []models.EquipmentNotification {
	var notifications []models.EquipmentNotification

	// Add equipment alerts
	for _, alert := range r.alerts {
		notifications = append(notifications, alert)
	}

	// Add maintenance reminders
	for _, reminder := range r.maintenanceReminders {
		notifications = append(notifications, reminder)
	}

	// Add status updates
	for _, statusUpdate := range r.statusUpdates {
		notifications = append(notifications, statusUpdate)
	}

	// Add performance alerts
	for _, perfAlert := range r.performanceAlerts {
		notifications = append(notifications, perfAlert)
	}

	return notifications
}

// processOrgChain recursively processes the organization chain input
func (r *Resolver) processOrgChain(facility *models.Facility, input models.AddOrgChainInput, parentManager *models.Manager) ([]models.Personnel, error) {
	var createdPersonnel []models.Personnel

	// Process managers at this level
	for _, managerInput := range input.Manager {
		manager, err := r.createManager(facility, *managerInput)
		if err != nil {
			return nil, err
		}
		createdPersonnel = append(createdPersonnel, manager)
	}

	// Process associates at this level
	for _, associateInput := range input.Associate {
		associate, err := r.createAssociate(facility, *associateInput, parentManager)
		if err != nil {
			return nil, err
		}
		createdPersonnel = append(createdPersonnel, associate)
	}

	// Process next level recursively
	if input.NextLevel != nil {
		// Find the most recently created manager to be the parent for the next level
		var nextLevelParent *models.Manager
		for i := len(createdPersonnel) - 1; i >= 0; i-- {
			if manager, ok := createdPersonnel[i].(*models.Manager); ok {
				nextLevelParent = manager
				break
			}
		}

		nextLevelPersonnel, err := r.processOrgChain(facility, *input.NextLevel, nextLevelParent)
		if err != nil {
			return nil, err
		}
		createdPersonnel = append(createdPersonnel, nextLevelPersonnel...)
	}

	return createdPersonnel, nil
}

// createManager creates a new manager and adds it to the facility
func (r *Resolver) createManager(facility *models.Facility, input models.AddManagerInput) (*models.Manager, error) {
	managerID := fmt.Sprintf("manager-%d", len(r.managers)+1)
	manager := &models.Manager{
		ID:            managerID,
		Name:          input.Name,
		Email:         input.Email,
		Phone:         input.Phone,
		JoinedAt:      input.JoinedAt,
		Status:        models.PersonnelStatusActive,
		Department:    input.Department,
		DirectReports: 0,
		Level:         input.Level,
	}

	// Add to resolver
	r.managers = append(r.managers, manager)
	r.personnel = append(r.personnel, manager)

	// Add to facility personnel
	facility.Personnel = append(facility.Personnel, manager)

	return manager, nil
}

// createAssociate creates a new associate and adds it to the facility
func (r *Resolver) createAssociate(facility *models.Facility, input models.AddAssociateInput, reportsTo *models.Manager) (*models.Associate, error) {
	associateID := fmt.Sprintf("associate-%d", len(r.associates)+1)
	associate := &models.Associate{
		ID:         associateID,
		Name:       input.Name,
		Email:      input.Email,
		Phone:      input.Phone,
		JoinedAt:   input.JoinedAt,
		Status:     models.PersonnelStatusActive,
		JobTitle:   input.JobTitle,
		Department: input.Department,
		ReportsTo:  reportsTo,
	}

	// Add to resolver
	r.associates = append(r.associates, associate)
	r.personnel = append(r.personnel, associate)

	// Add to facility personnel
	facility.Personnel = append(facility.Personnel, associate)

	// Update manager's direct reports count if reportsTo is specified
	if reportsTo != nil {
		reportsTo.DirectReports++
	}

	return associate, nil
}
