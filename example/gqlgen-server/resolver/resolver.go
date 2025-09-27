package resolver

import (
	"time"

	"github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	equipment          []*models.Equipment
	facilities         []*models.Facility
	maintenanceRecords []*models.MaintenanceRecord
	operationalMetrics []*models.OperationalMetric
	alerts             []*models.EquipmentAlert
}

// NewResolver creates a new resolver with sample data
func NewResolver() *Resolver {
	r := &Resolver{
		equipment:          make([]*models.Equipment, 0),
		facilities:         make([]*models.Facility, 0),
		maintenanceRecords: make([]*models.MaintenanceRecord, 0),
		operationalMetrics: make([]*models.OperationalMetric, 0),
		alerts:             make([]*models.EquipmentAlert, 0),
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
}
