package val

var property_info_sm_all = []string{"climb_stairs", "potential_for_noise", "pet_live_property", "park_on_property", "spaces_are_shared", "lack_of_electricity", "amenity_limit", "weapon_on_property"}

var safety_consider_sm_all = []string{"unsuitable_for_children_2_12", "unsuitable_for_infants_under_2", "pool_no_gate", "nearby_water", "climb_play_structure", "heights_no_rails", "dangerous_animal"}

var safety_devices_sm_all = []string{"cameras_audio_devices", "smoke_alarm", "carbon_monoxide_alarm"}

var smData = map[string][]string{"property_info": property_info_sm_all, "safety_consider": safety_consider_sm_all, "safety_devices": safety_devices_sm_all}

var sm_has_details = []string{"cameras_audio_devices","smoke_alarm","carbon_monoxide_alarm"}