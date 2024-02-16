package algo

var bedAndBreakfast = []string{"Pleasant", "Hidden", "Eclectic", "Delightful", "Hospitable", "Home away from home", "Relaxed", "Whimsical", "Country-style", "Personalized"}

var boats = []string{"Reliable", "Capable", "Sporty", "Yacht-like", "Seafaring", "Unsinkable", "Classic", "Cutting-edge", "Ocean-going", "Seaworthy", "Responsive", "State-of-the-art", "Streamlined", "Sturdy", "Nautical", "Swift", "Adventure-seeking", "Effortless", "Immaculate", "Stable", "Agile", "High-performance", "Compact", "Well-equipped", "Robust", "Adventure-ready", "Floating", "Seagoing", "Sail-powered", "Flexible", "Incredible", "Fast", "Advanced", "Speedy", "Sailing", "Aesthetically pleasing", "Versatile", "Powerful"}

var earthHomes = []string{"Eco-conscious", "Grounded", "Blending", "Simple", "Harmonious", "Earthen", "Wholesome", "Renewable", "Resilient", "Earthy", "Balanced", "Integrated", "Practical", "Humble", "Handcrafted", "Connected", "Organic", "Holistic"}

var beach = []string{"Seashell-studded", "Sun-kissed", "Windswept", "Remote", "Infinite", "Laid-back", "Golden", "Sun-drenched", "Lush", "Rhythmic", "Glistening", "Unspoiled", "Serenading", "Azure", "Breath-taking", "Crystal-clear", "Pure", "Untouched", "Ethereal", "White-sand", "Endless", "Sparkling", "Sandy", "Calm", "Playful", "Crisp", "Seashell-strewn", "Undisturbed", "Wondrous", "Sunlit"}

var historicalHomes = []string{"Antiquated", "Noble", "Fascinating", "Intricate", "Antique", "Ornate", "Significant", "Stately", "Classical", "Cherished", "Enduring", "Well-preserved", "Storied", "Regal", "Fabled", "Heritage", "Resplendent", "Old-world", "Time-honored", "Monumental", "Venerable", "Awe-inspiring", "Landmarked", "Glorious", "Preserved", "Epic", "Cultural", "Archaeological"}

var privateRooms = []string{"Personal", "Well-furnished", "Lodge", "Condo", "Private"}

var hotel = []string{"Five-star", "Impeccable", "Prestigious", "Glamorous", "Panoramic", "Resort-like"}

var unique = []string{"Iconic", "One-of-a-kind", "Unusual", "Futuristic", "Remarkable", "Rare", "Avant-garde", "Unrivaled", "Creative", "Unmatched", "Revolutionary", "Original", "Distinctive", "Uniquely designed", "Visionary", "Incomparable", "Eccentric", "Uncommon", "Distinct", "Outstanding", "Bold", "Unconventional", "Groundbreaking", "Extraordinary", "Unparalleled", "Unorthodox", "Unprecedented"}

var sharedRooms = []string{"Collaborative", "Roommate", "Boarding-house", "Cooperative", "Co-op", "Shared-living", "Bunkhouse", "House-share", "Shared-accommodation", "Communal", "Apartment-share", "Collective", "Multi-occupancy", "Hostel", "Joint", "Dormitory", "Shared", "Co-residence", "Co-housing", "Shared-residence", "Flatshare", "Shared-unit", "Group", "Co-living", "Co-tenancy", "Room-sharing", "Co-habitation", "Rooming-house", "Shared-dwelling", "Cohabitant", "Group-living", "Shared-space", "Shared-facility"}

var resort = []string{"Extravagant", "Haven", "Indulgent", "Leisure", "Well-being", "Stunning", "Vacation", "Unwind", "Recreational", "Wellness", "Getaway", "Sanctuary", "First-class", "Escape", "Pampering", "Oasis", "Spectacular", "All-inclusive", "Retreat", "Lavish", "Resort-style"}

var house = []string{"Abode", "Dwelling", "Flat", "Airy", "Urban", "Duplex", "Loft", "Pet-friendly", "Desirable", "Farmhouse", "Triplex", "Countryside", "Seaside", "Metropolitan", "Suburban", "Gated", "Convenient", "Estate", "Mountain", "Diverse", "Lakefront", "Manor", "Bungalow", "Homestead", "Waterfront", "Affordable", "Multicultural", "Studio", "Housing", "Village", "Terrace", "Beachfront", "Cottage", "Mansion", "Penthouse", "Accessible", "City", "Chalet", "Bright", "Renovated", "Home", "Rowhouse", "Townhouse", "Residence", "Ranch", "Family-friendly", "Child-friendly", "Well-connected", "Villa", "Apartment", "Country"}

var studio = []string{"Innovative", "Creative", "Collaborative", "Dynamic", "Inspiring", "State-of-the-art", "Cutting-edge", "Tech-savvy", "Futuristic", "Trendsetting", "Adaptable", "Versatile", "Resourceful", "Productive", "Organized", "Efficient", "Contemporary", "Stylish", "Aesthetic", "Ambitious", "Expressive", "Multifunctional", "Vibrant", "Energetic", "Spacious", "Ergonomic", "Comfortable", "Sleek", "Trendy", "Thoughtful", "Progressive", "Engaging", "Pioneering", "Imaginative", "Tech-driven", "Revolutionary", "Elegant", "Minimalist", "Artistic", "Invigorating", "High-tech", "Automated", "Smart", "Interactive", "Seamless", "Streamlined", "Unconventional", "Inspiring", "Empowering", "Modern", "Visionary", "User-friendly", "Groundbreaking", "Adaptable", "Digital", "Trendsetting", "Unprecedented", "Revolutionary", "Exceptional", "Inventive", "Progressive", "Customizable", "Evocative", "Elevated", "Fluid", "Contemporary", "Astute", "Cutting-edge", "Inclusive", "Strategic", "Futuristic", "Distinctive", "Evolving", "Augmented", "Crafted", "Refined", "Unmatched", "Adaptive", "Dynamic", "Resilient", "Exquisite", "Sustainable", "Distinctive", "Unparalleled", "Intuitive", "Captivating", "Enlightening", "Insightful", "Connected", "Multidimensional", "Hyper-modern", "Sculpted", "Streamlined", "Artful", "Visionary", "Elevated", "Strategic", "Ingenious", "Visionary"}

var optionDes = map[string][]string{"bed_and_breakfast": bedAndBreakfast, "boat": boats, "earth_home": earthHomes, "historical_home": historicalHomes, "private_room": privateRooms, "hotel": hotel, "unique": unique, "shared_room": sharedRooms, "resort": resort, "house": house, "beach": beach, "studio": studio}

// Amenities

var bedAndBreakfastAm = []string{"garden", "fire_pit", "hammock"}
var earthHomeAm = []string{"garden", "fire_pit", "hammock", "kayak", "lake_access", "arcade", "bikes", "books_and_reading_materials", "climbing_wall"}
var houseAm = []string{"baking_sheet", "barbecue_utensil", "garden", "bikes", "books_and_reading_materials"}
var boatAm = []string{"beach_essential", "boat_berth", "kayak"}
var beachAm = []string{"beach_essential", "boat_berth", "hammock", "beach_access", "resort_access", "sauna"}
var historicalHomeAm = []string{"sauna", "arcade_games", "piano", "cinema", "theme_room", "art", "music", "books_and_reading_materials", "climbing_wall", "bikes"}

var studioAm = []string{"sauna", "arcade_games", "piano", "cinema", "theme_room", "art", "music", "books_and_reading_materials"}

var privateRoomAm = []string{"theme_room"}
var hotelAm = []string{"sauna", "basketball_court", "tennis_court", "football_court"}
var uniqueAm = []string{"hammock", "lake_access", "sauna", "arcade_games", "bowling_alley", "batting_cage", "climbing_wall", "piano", "cinema"}
var sharedRoomAm = []string{"theme_room"}
var resortAm = []string{"beach_essential", "bikes", "hammock", "kayak", "beach_access", "resort_access", "sauna", "arcade_games", "bowling_alley", "batting_cage", "cinema", "theme_room", "pool", "tennis_court", "basketball_court", "football_court"}

var optionAmenity = map[string][]string{"bed_and_breakfast": bedAndBreakfastAm, "boat": boatAm, "earth_home": earthHomeAm, "historical_home": historicalHomeAm, "private_room": privateRoomAm, "hotel": hotelAm, "unique": uniqueAm, "shared_room": sharedRoomAm, "resort": resortAm, "house": houseAm, "beach": beachAm, "studio": studioAm}

// SpacesAreas

var resortSpace = []string{"pool", "hot_tub", "gym", "bedroom", "bathroom"}
var hotelSpace = []string{"gym", "pool", "hot_tub", "bedroom", "bathroom"}
var historicalHomeSpace = []string{"office", "garden", "bedroom", "bathroom"}
var bedAndBreakfastSpace = []string{"bedroom", "bathroom"}
var houseSpace = []string{"bedroom", "bathroom", "kitchen", "living_room", "dining_area"}
var boatSpace = []string{"bedroom"}

var studioSpace = []string{"office", "bedroom"}

var beachSpace = []string{"bathroom", "bathroom", "pool"}
var privateRoomSpace = []string{"bedroom", "bathroom"}
var uniqueSpace = []string{"bedroom"}
var sharedRoomSpace = []string{"bedroom"}
var earthHomeSpace = []string{"bedroom", "bathroom", "kitchen"}

var optionSpaceArea = map[string][]string{"bed_and_breakfast": bedAndBreakfastSpace, "boat": boatSpace, "earth_home": earthHomeSpace, "historical_home": historicalHomeSpace, "private_room": privateRoomSpace, "hotel": hotelSpace, "unique": uniqueSpace, "shared_room": sharedRoomSpace, "resort": resortSpace, "house": houseSpace, "beach": beachSpace, "studio": studioSpace}

// Highlights
var bedAndBreakfastHigh = []string{"family_first", "peaceful"}
var earthHomeHigh = []string{"artistic", "adventurous", "historical", "popular_location"}
var houseHigh = []string{"peaceful", "family_first", "bold", "eye_catching", "tech_first", "spacious"}
var boatHigh = []string{"adventurous", "bold"}
var beachHigh = []string{"popular_location", "eye_catching"}

var studioHigh = []string{"popular_location", "eye_catching", "tech_first", "spacious"}

var historicalHomeHigh = []string{"artistic", "historical", "eye_catching", "bold", "popular_location"}
var privateRoomHigh = []string{"peaceful", "spacious", "popular_location"}
var hotelHigh = []string{"tech_like", "eye_catching", "bold", "artistic", "popular_location", "historical", "family_first"}
var uniqueHigh = []string{"bold", "artistic", "adventurous", "eye_catching", "tech_like"}
var sharedRoomHigh = []string{"peaceful", "popular_location"}
var resortHigh = []string{"spacious", "tech_like", "eye_catching", "bold", "artistic", "popular_location", "historical", "family_first"}

var optionHighlight = map[string][]string{"bed_and_breakfast": bedAndBreakfastHigh, "boat": boatHigh, "earth_home": earthHomeHigh, "historical_home": historicalHomeHigh, "private_room": privateRoomHigh, "hotel": hotelHigh, "unique": uniqueHigh, "shared_room": sharedRoomHigh, "resort": resortHigh, "house": houseHigh, "beach": beachHigh, "studio": studioHigh}

// Shortlet type
var bedAndBreakfastType = []string{"bed_and_breakfast", "guest_house", "nature_lodge", "minsu", "casa"}

var studioType = []string{"studio"}

var earthHomeType = []string{"barn", "hut", "tree_house", "cabin"}
var houseType = []string{"flat_apartment", "home", "bungalow"}
var boatType = []string{"yacht"}
var beachType = []string{"nature_lodge", "cottage", "beach_house", "house_boat"}
var historicalHomeType = []string{"tower", "castle", "windmill", "lighthouse", "dome"}
var privateRoomType = []string{"flat_apartment", "guest_house", "loft", "equipped_apartment", "home", "bungalow"}
var hotelType = []string{"equipped_apartment", "hotel_room"}
var uniqueType = []string{"barn", "tiny_home", "plane", "tower", "bus", "container", "cave", "tree_house", "windmill", "motorhome"}
var sharedRoomType = []string{"flat_apartment", "guest_house", "loft", "equipped_apartment", "home", "bungalow"}
var resortType = []string{"home", "minsu", "hotel_room", "resort", "beach_house", "guest_house", "castle", "bed_and_breakfast"}

var optionType = map[string][]string{"bed_and_breakfast": bedAndBreakfastType, "boat": boatType, "earth_home": earthHomeType, "historical_home": historicalHomeType, "private_room": privateRoomType, "hotel": hotelType, "unique": uniqueType, "shared_room": sharedRoomType, "resort": resortType, "house": houseType, "beach": beachType, "studio": studioType}

// SpaceType
var bedAndBreakfastSpaceType = []string{"private_room", "shared_room", "full_place"}
var earthHomeSpaceType = []string{"private_room", "shared_room", "full_place"}
var houseSpaceType = []string{"full_place"}
var boatSpaceType = []string{"private_room", "shared_room", "full_place"}

var studioSpaceType = []string{"private_room", "shared_room", "full_place"}

var beachSpaceType = []string{"private_room", "shared_room", "full_place"}
var historicalHomeSpaceType = []string{"private_room", "shared_room", "full_place"}
var privateRoomSpaceType = []string{"private_room"}
var hotelSpaceType = []string{"private_room", "shared_room", "full_place"}
var uniqueSpaceType = []string{"private_room", "shared_room", "full_place"}
var sharedRoomSpaceType = []string{"shared_room"}
var resortSpaceType = []string{"private_room", "shared_room", "full_place"}

var optionSpaceType = map[string][]string{"bed_and_breakfast": bedAndBreakfastSpaceType, "boat": boatSpaceType, "earth_home": earthHomeSpaceType, "historical_home": historicalHomeSpaceType, "private_room": privateRoomSpaceType, "hotel": hotelSpaceType, "unique": uniqueSpaceType, "shared_room": sharedRoomSpaceType, "resort": resortSpaceType, "house": houseSpaceType, "beach": beachSpaceType, "studio": studioSpaceType}

var OptionCategory = []string{"house", "bed_and_breakfast", "boat", "unique", "resort", "historical_homes", "private_room", "earth_homes", "shared_room", "beach", "studio", "hotel"}
