package val


var host_event_cancel_reasons_all = []string{"low_ticket_sales","vendor_availability","venue_issues","force_majeure","health_safety_concerns","financial_constraints","logistical_challenges","government_regulations","conflict_of_interest","insufficient_planning_time","technical_difficulties","participant_safety","organizer_health_issues","unsatisfactory_event_details","sponsorship_withdrawal","legal_issues","other"}


var host_option_cancel_reasons_all = []string{"emergency_repairs","double_booking","health_and_safety_concerns","property_damage","personal_reasons","booking_violations","extenuating_circumstances","unexpected_property_sale","legal_or_regulatory_issues","payment_issues","guest_violations","unavailability","health_and_safety_standards","cancellation_policy","maintenance_work","miscommunication","security_concerns","other"}

var user_event_cancel_reasons_all = []string{"scheduling_conflict","medical_emergency","travel_issues","family_obligations","work_commitments","financial_constraints","personal_health_issues","event_rescheduling","transportation_problems","unforeseen_circumstances","event_changes","personal_reasons","safety_concerns","weather_conditions","change_in_plans","conflict_of_interest","ticket_purchase_error","unsatisfactory_event_details","event_cancellation","other"}


var user_option_cancel_reasons_all = []string{"change_in_plans","emergencies","illness_or_health_issues","travel_restrictions","financial_constraints","double_booking","weather_conditions","family_matters","work_commitments","fear_or_anxiety","transportation_issues","travel_document_issues","event_cancellations","health_and_safety","change_in_group_size","personal_reasons","discounts_or_better_offers","miscommunication","cancellation_policies","other"}

func ValidateHostEventCancelTypes(r string) bool {
	for _, a := range host_event_cancel_reasons_all {
		if r == a {
			return true
		}
	}
	return false
}


func ValidateUserEventCancelTypes(r string) bool {
	for _, a := range user_event_cancel_reasons_all {
		if r == a {
			return true
		}
	}
	return false
}

func ValidateHostOptionCancelTypes(r string) bool {
	for _, a := range host_option_cancel_reasons_all {
		if r == a {
			return true
		}
	}
	return false
}

func ValidateUserOptionCancelTypes(r string) bool {
	for _, a := range user_option_cancel_reasons_all {
		if r == a {
			return true
		}
	}
	return false
}