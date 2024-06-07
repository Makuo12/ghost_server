CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "cube";
CREATE EXTENSION IF NOT EXISTS "earthdistance";


CREATE TABLE "users" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "firebase_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "public_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "hashed_password" varchar NOT NULL,
  "deep_link_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "firebase_password" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone_number" varchar NOT NULL DEFAULT 'none',
  "first_name" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "last_name" varchar NOT NULL,
  "date_of_birth" date NOT NULL,
  "dial_code" varchar NOT NULL DEFAULT 'none',
  "dial_country" varchar NOT NULL DEFAULT 'none',
  "current_option_id" varchar NOT NULL DEFAULT 'none',
  "currency" varchar NOT NULL,
  "default_card" varchar NOT NULL DEFAULT 'none',
  "default_payout_card" varchar NOT NULL DEFAULT 'none',
  "default_account_id" varchar NOT NULL DEFAULT 'none',
  "is_active" boolean NOT NULL DEFAULT true,
  "is_deleted" boolean NOT NULL DEFAULT false,
  "photo" varchar NOT NULL DEFAULT 'none',
  "password_changed_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_apn_details" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "device_name" varchar NOT NULL,
  "model" varchar NOT NULL,
  "identifier_for_vendor" varchar NOT NULL,
  "token" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users_profiles" (
  "user_id" uuid PRIMARY KEY NOT NULL,
  "work" varchar NOT NULL DEFAULT 'none',
  "languages" varchar[] NOT NULL,
  "bio" text NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users_options_reviews" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_user_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "message" text NOT NULL DEFAULT 'none',
  "clean" decimal(4,2) NOT NULL DEFAULT '0.00',
  "communication" decimal(4,2) NOT NULL DEFAULT '0.00',
  "house_rules" decimal(4,2) NOT NULL DEFAULT '0.00',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "identity" (
  "user_id" uuid PRIMARY KEY NOT NULL,
  "country" varchar NOT NULL DEFAULT 'none',
  "type" varchar NOT NULL DEFAULT 'none',
  "id_photo" varchar NOT NULL DEFAULT 'none',
  "id_photo_list" varchar[] NOT NULL,
  "id_back_photo" varchar NOT NULL DEFAULT 'none',
  "id_back_photo_list" varchar[] NOT NULL,
  "facial_photo" varchar NOT NULL DEFAULT 'none',
  "facial_photo_list" varchar[] NOT NULL,
  "status" varchar NOT NULL DEFAULT 'not_started',
  "is_verified" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "em_contacts" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "name" varchar NOT NULL,
  "relationship" varchar NOT NULL,
  "email" varchar NOT NULL DEFAULT 'none',
  "dial_code" varchar NOT NULL DEFAULT 'none',
  "dial_country" varchar NOT NULL DEFAULT 'none',
  "phone_number" varchar NOT NULL DEFAULT 'none',
  "language" varchar NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "mailing_addresses" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "street" text NOT NULL,
  "city" varchar NOT NULL,
  "state" varchar NOT NULL,
  "country" varchar NOT NULL,
  "postcode" varchar NOT NULL,
  "geolocation" point NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "account_numbers" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "account_number" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "account_name" varchar NOT NULL,
  "recipient_code" varchar NOT NULL,
  "bank_name" varchar NOT NULL,
  "bank_code" varchar NOT NULL,
  "country" varchar NOT NULL,
  "type" varchar NOT NULL,
  "bank_id" int NOT NULL
);

CREATE TABLE "payments_gate_pays" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "transaction_id" bigint UNIQUE NOT NULL,
  "reference" varchar UNIQUE NOT NULL,
  "requested_amount" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "payment_gate_fee" bigint NOT NULL,
  "was_refunded" boolean NOT NULL DEFAULT false,
  "authorization_code" varchar UNIQUE NOT NULL,
  "payment_gate_paid_at" timestamptz NOT NULL,
  "channel" varchar NOT NULL,
  "payment_gate_created_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "cards" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "email" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "authorization_code" varchar UNIQUE NOT NULL,
  "card_type" varchar NOT NULL,
  "last4" varchar NOT NULL,
  "exp_month" varchar NOT NULL,
  "exp_year" varchar NOT NULL,
  "bank" varchar NOT NULL,
  "country_code" varchar NOT NULL,
  "reusable" boolean NOT NULL,
  "channel" varchar NOT NULL,
  "card_signature" varchar NOT NULL,
  "account_name" varchar NOT NULL,
  "bin" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "accounts" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "id_int" bigserial UNIQUE NOT NULL,
  "user_id" uuid NOT NULL,
  "currency" varchar NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users_locations" (
  "user_id" uuid PRIMARY KEY NOT NULL,
  "street" text NOT NULL,
  "city" varchar NOT NULL,
  "state" varchar NOT NULL,
  "country" varchar NOT NULL,
  "postcode" varchar NOT NULL,
  "geolocation" point NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "account_id" uuid NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "from_account_id" uuid NOT NULL,
  "to_account_id" uuid NOT NULL,
  "from_account_id_int" bigint NOT NULL,
  "to_account_id_int" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_questions" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "host_as_individual" boolean NOT NULL,
  "organization_name" varchar NOT NULL DEFAULT 'none',
  "organization_email" varchar NOT NULL DEFAULT 'none',
  "legal_represents" varchar[] NOT NULL,
  "street" varchar NOT NULL DEFAULT 'none',
  "city" varchar NOT NULL DEFAULT 'none',
  "state" varchar NOT NULL DEFAULT 'none',
  "country" varchar NOT NULL DEFAULT 'none',
  "postcode" varchar NOT NULL DEFAULT 'none',
  "geolocation" point NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "wishlists" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "wishlists_items" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "wishlist_id" uuid NOT NULL,
  "option_user_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "vids" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "path" varchar NOT NULL,
  "filter" varchar NOT NULL,
  "option_user_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "main_option_type" varchar NOT NULL,
  "start_date" varchar NOT NULL,
  "caption" varchar NOT NULL,
  "from_who" varchar NOT NULL,
  "extra_option_id" uuid NOT NULL,
  "extra_option_id_fake" boolean NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "options_infos" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "co_host_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "option_user_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "host_id" uuid NOT NULL,
  "deep_link_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "primary_user_id" uuid NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_complete" boolean NOT NULL DEFAULT false,
  "is_verified" boolean NOT NULL DEFAULT true,
  "category" varchar NOT NULL DEFAULT 'none',
  "category_two" varchar NOT NULL DEFAULT 'none',
  "category_three" varchar NOT NULL DEFAULT 'none',
  "category_four" varchar NOT NULL DEFAULT 'none',
  "is_top_seller" boolean NOT NULL DEFAULT false,
  "time_zone" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "option_img" varchar NOT NULL,
  "option_type" varchar NOT NULL,
  "main_option_type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "completed" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "options_infos_category" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "type_of_shortlet" varchar[] NOT NULL,
  "amenities" varchar[] NOT NULL,
  "highlight" varchar[] NOT NULL,
  "space_area" varchar[] NOT NULL,
  "space_type" varchar[] NOT NULL,
  "des" varchar[] NOT NULL,
  "name" varchar[] NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "events_infos_category" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "event_type" varchar[] NOT NULL,
  "event_sub_type" varchar[] NOT NULL,
  "highlight" varchar[] NOT NULL,
  "des" varchar[] NOT NULL,
  "name" varchar[] NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "options_infos_status" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "status" varchar NOT NULL DEFAULT 'unlist',
  "status_reason" varchar NOT NULL DEFAULT 'none',
  "snooze_start_date" date NOT NULL,
  "snooze_end_date" date NOT NULL,
  "unlist_reason" varchar NOT NULL DEFAULT 'none',
  "unlist_des" varchar NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "complete_option_info" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "current_state" varchar NOT NULL DEFAULT 'none',
  "previous_state" varchar NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "options_info_details" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "des" text NOT NULL DEFAULT 'none',
  "space_des" text NOT NULL DEFAULT 'none',
  "guest_access_des" text NOT NULL DEFAULT 'none',
  "interact_with_guests_des" text NOT NULL DEFAULT 'none',
  "pets_allowed" boolean NOT NULL DEFAULT false,
  "other_des" text NOT NULL DEFAULT 'none',
  "neighborhood_des" text NOT NULL DEFAULT 'none',
  "get_around_des" text NOT NULL DEFAULT 'none',
  "host_name_option" varchar NOT NULL DEFAULT 'none',
  "option_highlight" varchar[] NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_add_charges" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "type" varchar NOT NULL,
  "main_fee" bigint NOT NULL,
  "extra_fee" bigint NOT NULL,
  "num_of_guest" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_discounts" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "type" varchar NOT NULL,
  "main_type" varchar NOT NULL,
  "percent" int NOT NULL,
  "name" varchar NOT NULL DEFAULT 'none',
  "extra_type" varchar NOT NULL DEFAULT 'none',
  "des" text NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "options_info_photos" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "cover_image" varchar NOT NULL,
  "has_meta_data" boolean NOT NULL DEFAULT false,
  "public_cover_image" varchar NOT NULL DEFAULT 'none',
  "public_photo" varchar[] NOT NULL,
  "photo" varchar[] NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "options_photo_captions" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "photo_id" varchar UNIQUE NOT NULL,
  "caption" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_infos" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "sub_category_type" varchar NOT NULL DEFAULT 'none',
  "event_type" varchar NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_date_times" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "event_info_id" uuid NOT NULL,
  "start_date" date NOT NULL,
  "name" varchar NOT NULL DEFAULT 'none',
  "publish_check_in_steps" boolean NOT NULL DEFAULT false,
  "check_in_method" varchar NOT NULL DEFAULT 'none',
  "event_dates" varchar[] NOT NULL,
  "deep_link_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "type" varchar NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "need_bands" boolean NOT NULL DEFAULT false,
  "need_tickets" boolean NOT NULL DEFAULT true,
  "absorb_band_charge" boolean NOT NULL DEFAULT false,
  "status" varchar NOT NULL DEFAULT 'on_sale',
  "note" text NOT NULL DEFAULT 'none',
  "end_date" date NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_check_in_steps" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "event_date_time_id" uuid NOT NULL,
  "photo" varchar NOT NULL DEFAULT 'none',
  "des" text NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_date_details" (
  "event_date_time_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "start_time" varchar NOT NULL DEFAULT 'none',
  "end_time" varchar NOT NULL DEFAULT 'none',
  "time_zone" varchar NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_date_publishes" (
  "event_date_time_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "event_public" varchar NOT NULL DEFAULT 'public',
  "event_going_public" varchar NOT NULL DEFAULT 'no',
  "event_going_public_date" date NOT NULL,
  "event_going_public_time" time NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_date_private_audiences" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "event_date_time_id" uuid NOT NULL,
  "name" varchar NOT NULL,
  "type" varchar NOT NULL,
  "email" varchar NOT NULL DEFAULT 'none',
  "number" varchar NOT NULL DEFAULT 'none',
  "sent" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_date_tickets" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "event_date_time_id" uuid NOT NULL,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "start_time" time NOT NULL,
  "end_time" time NOT NULL,
  "deep_link_id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "price" bigint NOT NULL,
  "absorb_fees" boolean NOT NULL,
  "description" text NOT NULL,
  "capacity" int NOT NULL,
  "capacity_sold" int NOT NULL DEFAULT 0,
  "type" varchar NOT NULL,
  "level" varchar NOT NULL,
  "ticket_type" varchar NOT NULL,
  "num_of_seats" int NOT NULL DEFAULT 0,
  "free_refreshment" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_date_locations" (
  "event_date_time_id" uuid PRIMARY KEY NOT NULL,
  "street" text NOT NULL,
  "city" varchar NOT NULL,
  "state" varchar NOT NULL,
  "country" varchar NOT NULL,
  "postcode" varchar NOT NULL,
  "geolocation" point NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "locations" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "street" varchar NOT NULL,
  "city" varchar NOT NULL,
  "state" varchar NOT NULL,
  "country" varchar NOT NULL,
  "postcode" varchar NOT NULL,
  "geolocation" point NOT NULL,
  "show_specific_location" boolean NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "shortlets" (
  "option_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "space_type" varchar NOT NULL DEFAULT 'none',
  "any_space_shared" boolean NOT NULL DEFAULT false,
  "guest_welcomed" int NOT NULL,
  "publish_check_in_steps" boolean NOT NULL DEFAULT false,
  "year_built" int NOT NULL,
  "check_in_method" varchar NOT NULL DEFAULT 'code_scan',
  "check_in_method_des" text NOT NULL DEFAULT 'none',
  "property_size" int NOT NULL,
  "shared_spaces_with" varchar[] NOT NULL,
  "property_size_unit" varchar NOT NULL,
  "type_of_shortlet" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_date_times" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "date" date NOT NULL,
  "available" boolean NOT NULL,
  "price" bigint NOT NULL,
  "note" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "space_areas" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "shared_space" boolean NOT NULL DEFAULT false,
  "space_type" varchar NOT NULL,
  "photos" varchar[] NOT NULL,
  "beds" varchar[] NOT NULL,
  "is_suite" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "options_prices" (
  "option_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "price" bigint NOT NULL,
  "weekend_price" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "wifi_details" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "network_name" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "amenities" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "tag" varchar NOT NULL,
  "am_type" varchar NOT NULL,
  "has_am" boolean NOT NULL DEFAULT false,
  "time_set" boolean NOT NULL DEFAULT false,
  "location_option" varchar NOT NULL DEFAULT 'none',
  "size_option" int NOT NULL DEFAULT 0,
  "privacy_option" varchar NOT NULL DEFAULT 'none',
  "time_option" varchar NOT NULL DEFAULT 'none',
  "start_time" time NOT NULL DEFAULT (now()),
  "end_time" time NOT NULL DEFAULT (now()),
  "availability_option" varchar NOT NULL DEFAULT 'none',
  "start_month" varchar NOT NULL DEFAULT 'none',
  "end_month" varchar NOT NULL DEFAULT 'none',
  "type_option" varchar NOT NULL DEFAULT 'none',
  "price_option" varchar NOT NULL DEFAULT 'none',
  "brand_option" varchar NOT NULL DEFAULT 'none',
  "list_options" varchar[] NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "check_in_steps" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "photo" varchar NOT NULL DEFAULT 'none',
  "des" text NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_availability_settings" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "advance_notice" varchar NOT NULL DEFAULT 'same day',
  "auto_block_dates" boolean NOT NULL DEFAULT true,
  "advance_notice_condition" varchar NOT NULL DEFAULT 'any time',
  "preparation_time" varchar NOT NULL DEFAULT 'none',
  "availability_window" varchar NOT NULL DEFAULT '12 months in advance',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "cancel_policies" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "type_one" varchar NOT NULL DEFAULT 'flexible',
  "type_two" varchar NOT NULL DEFAULT 'none',
  "request_a_refund" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_trip_lengths" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "min_stay_day" int NOT NULL DEFAULT 1,
  "max_stay_night" int NOT NULL DEFAULT 365,
  "manual_approve_request_pass_max" boolean NOT NULL DEFAULT true,
  "allow_reservation_request" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "check_in_out_details" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "arrive_after" varchar NOT NULL DEFAULT '10:30',
  "arrive_before" varchar NOT NULL DEFAULT '20:30',
  "leave_before" varchar NOT NULL DEFAULT '18:30',
  "restricted_check_in_days" varchar[] NOT NULL,
  "restricted_check_out_days" varchar[] NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_co_hosts" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "co_user_id" varchar NOT NULL DEFAULT 'none',
  "accepted" boolean NOT NULL DEFAULT false,
  "email" varchar NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "reservations" boolean NOT NULL DEFAULT false,
  "post" boolean NOT NULL DEFAULT false,
  "scan_code" boolean NOT NULL DEFAULT false,
  "calender" boolean NOT NULL DEFAULT false,
  "insights" boolean NOT NULL DEFAULT false,
  "edit_option_info" boolean NOT NULL DEFAULT false,
  "edit_event_dates_times" boolean NOT NULL DEFAULT false,
  "edit_co_hosts" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "options_extra_infos" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "type" varchar NOT NULL,
  "info" text NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "things_to_notes" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "checked" boolean NOT NULL,
  "tag" varchar NOT NULL,
  "type" varchar NOT NULL,
  "des" text NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_invitations" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "message" text NOT NULL DEFAULT 'none',
  "transport_method" varchar NOT NULL,
  "user_email" varchar NOT NULL,
  "user_phone_number" varchar NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_rules" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "checked" boolean NOT NULL,
  "tag" varchar NOT NULL,
  "type" varchar NOT NULL,
  "des" text NOT NULL DEFAULT 'none',
  "start_time" time NOT NULL,
  "end_time" time NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_book_methods" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "instant_book" boolean NOT NULL DEFAULT true,
  "identity_verified" boolean NOT NULL DEFAULT false,
  "good_track_record" boolean NOT NULL DEFAULT false,
  "pre_book_msg" text NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "book_requirements" (
  "option_id" uuid PRIMARY KEY NOT NULL,
  "email" boolean NOT NULL DEFAULT true,
  "phone_number" boolean NOT NULL DEFAULT true,
  "rules" boolean NOT NULL DEFAULT true,
  "payment_info" boolean NOT NULL DEFAULT true,
  "profile_photo" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_messages" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "message" text NOT NULL,
  "seen" boolean NOT NULL DEFAULT false,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "feedbacks" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "subject" varchar NOT NULL,
  "sub_subject" varchar NOT NULL DEFAULT 'none',
  "detail" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "helps" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "email" varchar NOT NULL,
  "subject" varchar NOT NULL,
  "sub_subject" varchar NOT NULL DEFAULT 'none',
  "detail" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "report_options" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "option_user_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "type_one" varchar NOT NULL DEFAULT 'none',
  "type_two" varchar NOT NULL DEFAULT 'none',
  "type_three" varchar NOT NULL DEFAULT 'none',
  "description" text NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "charge_references" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "reference" uuid NOT NULL DEFAULT (uuid_generate_v4()),
  "reason" varchar NOT NULL,
  "is_complete" boolean NOT NULL DEFAULT false,
  "charge" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "cancellations" (
  "charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "type" varchar NOT NULL,
  "main_option_type" varchar NOT NULL,
  "charge_type" varchar NOT NULL,
  "cancel_user_id" uuid NOT NULL,
  "reason_one" varchar NOT NULL,
  "reason_two" varchar NOT NULL DEFAULT 'none',
  "message" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "main_refunds" (
  "charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "user_percent" int NOT NULL,
  "host_percent" int NOT NULL,
  "charge_type" varchar NOT NULL,
  "type" varchar NOT NULL,
  "is_payed" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "refunds" (
  "charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "reference" varchar NOT NULL,
  "send_medium" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "user_id" uuid NOT NULL,
  "amount_payed" bigint NOT NULL,
  "time_paid" timestamptz NOT NULL DEFAULT (now()),
  "is_complete" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "refund_payouts" (
  "charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "amount" bigint NOT NULL,
  "service_fee" bigint NOT NULL,
  "user_id" uuid NOT NULL,
  "time_paid" timestamptz NOT NULL DEFAULT (now()),
  "currency" varchar NOT NULL,
  "account_number" varchar NOT NULL DEFAULT 'none',
  "is_complete" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "main_payouts" (
  "charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "type" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "time_paid" timestamptz NOT NULL DEFAULT (now()),
  "service_fee" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "account_number" varchar NOT NULL DEFAULT 'none',
  "is_complete" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payouts" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "payout_ids" uuid[] NOT NULL,
  "user_id" uuid NOT NULL,
  "send_medium" varchar NOT NULL,
  "parent_type" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "amount_payed" bigint NOT NULL,
  "account_number" varchar NOT NULL,
  "time_paid" timestamptz NOT NULL DEFAULT (now()),
  "transfer_code" varchar NOT NULL DEFAULT 'none',
  "is_complete" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "charge_option_references" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "option_user_id" uuid NOT NULL,
  "discount" varchar NOT NULL,
  "main_price" bigint NOT NULL,
  "service_fee" bigint NOT NULL,
  "total_fee" bigint NOT NULL,
  "date_price" varchar[] NOT NULL,
  "guests" varchar[] NOT NULL,
  "date_booked" timestamptz NOT NULL,
  "currency" varchar NOT NULL,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "guest_fee" bigint NOT NULL,
  "pet_fee" bigint NOT NULL,
  "clean_fee" bigint NOT NULL,
  "nightly_pet_fee" bigint NOT NULL,
  "nightly_guest_fee" bigint NOT NULL,
  "can_instant_book" boolean NOT NULL,
  "require_request" boolean NOT NULL,
  "request_type" varchar NOT NULL,
  "reference" varchar NOT NULL,
  "payment_reference" varchar NOT NULL,
  "request_approved" boolean NOT NULL,
  "is_complete" boolean NOT NULL,
  "cancelled" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "charge_event_references" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "option_user_id" uuid NOT NULL,
  "total_fee" bigint NOT NULL,
  "service_fee" bigint NOT NULL,
  "total_absorb_fee" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "date_booked" timestamptz NOT NULL,
  "can_instant_book" boolean NOT NULL DEFAULT true,
  "require_request" boolean NOT NULL DEFAULT false,
  "request_type" varchar NOT NULL DEFAULT 'none',
  "reference" varchar NOT NULL,
  "payment_reference" varchar NOT NULL,
  "request_approved" boolean NOT NULL DEFAULT false,
  "is_complete" boolean NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "charge_date_references" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "charge_event_id" uuid NOT NULL,
  "event_date_id" uuid NOT NULL,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "date_booked" timestamptz NOT NULL,
  "start_time" varchar NOT NULL,
  "end_time" varchar NOT NULL,
  "total_date_fee" bigint NOT NULL,
  "total_date_service_fee" bigint NOT NULL,
  "total_date_absorb_fee" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "charge_ticket_references" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "charge_date_id" uuid NOT NULL,
  "ticket_id" uuid NOT NULL,
  "grade" varchar NOT NULL,
  "type" varchar NOT NULL,
  "date_booked" timestamptz NOT NULL,
  "price" bigint NOT NULL,
  "service_fee" bigint NOT NULL,
  "absorb_fee" bigint NOT NULL,
  "ticket_type" varchar NOT NULL,
  "group_price" bigint NOT NULL,
  "gifted" varchar NOT NULL DEFAULT 'none',
  "cancelled" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "charge_reviews" (
  "charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "type" varchar NOT NULL,
  "general" int NOT NULL,
  "environment" int NOT NULL DEFAULT 0,
  "accuracy" int NOT NULL DEFAULT 0,
  "check_in" int NOT NULL DEFAULT 0,
  "communication" int NOT NULL DEFAULT 0,
  "is_published" boolean NOT NULL DEFAULT false,
  "location" int NOT NULL DEFAULT 0,
  "current_state" varchar NOT NULL DEFAULT 'none',
  "previous_state" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'started',
  "private_note" text NOT NULL DEFAULT 'none',
  "public_note" text NOT NULL DEFAULT 'none',
  "stay_clean" varchar NOT NULL DEFAULT 'none',
  "stay_comfort" varchar NOT NULL DEFAULT 'none',
  "host_review" varchar NOT NULL DEFAULT 'none',
  "amenities" varchar[] NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "complete_charge_reviews" (
  "charge_review_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "current_state" varchar NOT NULL DEFAULT 'none',
  "previous_state" varchar NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "scanned_charges" (
  "charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "scanned" boolean NOT NULL DEFAULT false,
  "scanned_by" uuid NOT NULL,
  "scanned_time" timestamptz NOT NULL,
  "charge_type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "option_reference_infos" (
  "option_charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "amenities" varchar[] NOT NULL,
  "space_area" varchar NOT NULL,
  "time_zone" varchar NOT NULL,
  "arrive_before" varchar NOT NULL,
  "arrive_after" varchar NOT NULL,
  "leave_before" varchar NOT NULL,
  "cancel_policy_one" varchar NOT NULL,
  "cancel_policy_two" varchar NOT NULL,
  "pets_allowed" boolean NOT NULL,
  "rules_checked" varchar[] NOT NULL,
  "rules_unchecked" varchar[] NOT NULL,
  "shortlet" varchar NOT NULL,
  "location" varchar NOT NULL,
  "host_as_individual" boolean NOT NULL,
  "organization_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "event_reference_infos" (
  "event_charge_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "event_date_location" varchar NOT NULL,
  "event_info" varchar NOT NULL,
  "event_date_times" varchar NOT NULL,
  "cancel_policy_one" varchar NOT NULL,
  "cancel_policy_two" varchar NOT NULL,
  "event_date_details" varchar NOT NULL,
  "host_as_individual" boolean NOT NULL,
  "organization_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "messages" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "msg_id" uuid UNIQUE NOT NULL,
  "sender_id" uuid NOT NULL,
  "receiver_id" uuid NOT NULL,
  "message" text NOT NULL DEFAULT 'none',
  "type" varchar NOT NULL,
  "read" boolean NOT NULL DEFAULT false,
  "photo" varchar NOT NULL DEFAULT 'none',
  "parent_id" varchar NOT NULL DEFAULT 'none',
  "reference" varchar NOT NULL DEFAULT 'none',
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "request_notifies" (
  "m_id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "start_date" varchar NOT NULL,
  "end_date" varchar NOT NULL,
  "has_price" boolean NOT NULL,
  "same_price" boolean NOT NULL,
  "price" bigint NOT NULL,
  "item_id" varchar NOT NULL,
  "approved" boolean NOT NULL DEFAULT false,
  "cancelled" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "single_rooms" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_one" uuid NOT NULL,
  "user_two" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "notifications" (
  "id" uuid UNIQUE PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "header" varchar NOT NULL,
  "item_id" uuid NOT NULL,
  "item_id_fake" bool NOT NULL,
  "user_id" uuid NOT NULL,
  "type" varchar NOT NULL,
  "message" text NOT NULL,
  "handled" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("id");

CREATE INDEX ON "users" ("last_name");

CREATE INDEX ON "users" ("first_name");

CREATE INDEX ON "users" ("is_active");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users_profiles" ("user_id");

CREATE INDEX ON "mailing_addresses" ("id");

CREATE INDEX ON "mailing_addresses" ("user_id");

CREATE INDEX ON "account_numbers" ("id");

CREATE INDEX ON "account_numbers" ("user_id");

CREATE INDEX ON "cards" ("id");

CREATE INDEX ON "cards" ("user_id");

CREATE INDEX ON "accounts" ("id");

CREATE INDEX ON "accounts" ("user_id");

CREATE INDEX ON "accounts" ("balance");

CREATE INDEX ON "users_locations" ("user_id");

CREATE INDEX ON "users_locations" ("city");

CREATE INDEX ON "users_locations" ("state");

CREATE INDEX ON "users_locations" ("country");

CREATE INDEX ON "entries" ("id");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "wishlists_items" ("id");

CREATE INDEX ON "wishlists_items" ("option_user_id");

CREATE INDEX ON "wishlists_items" ("created_at");

CREATE INDEX ON "options_infos" ("id");

CREATE INDEX ON "options_infos" ("host_id");

CREATE INDEX ON "options_infos" ("is_active");

CREATE INDEX ON "options_infos" ("is_verified");

CREATE INDEX ON "options_infos_category" ("option_id");

CREATE INDEX ON "events_infos_category" ("option_id");

CREATE INDEX ON "options_infos_status" ("option_id");

CREATE INDEX ON "complete_option_info" ("option_id");

CREATE INDEX ON "complete_option_info" ("created_at");

CREATE INDEX ON "options_info_details" ("option_id");

CREATE INDEX ON "options_info_details" ("host_name_option");

CREATE INDEX ON "options_info_details" ("option_highlight");

CREATE INDEX ON "options_info_details" ("created_at");

CREATE INDEX ON "options_info_photos" ("option_id");

CREATE INDEX ON "options_info_photos" ("created_at");

CREATE INDEX ON "options_photo_captions" ("id");

CREATE INDEX ON "options_photo_captions" ("option_id");

CREATE INDEX ON "options_photo_captions" ("created_at");

CREATE INDEX ON "event_infos" ("event_type");

CREATE INDEX ON "event_infos" ("sub_category_type");

CREATE INDEX ON "event_infos" ("option_id");

CREATE INDEX ON "event_infos" ("created_at");

CREATE INDEX ON "event_date_times" ("event_info_id");

CREATE INDEX ON "event_date_times" ("id");

CREATE INDEX ON "event_date_times" ("created_at");

CREATE INDEX ON "event_date_times" ("start_date");

CREATE INDEX ON "event_date_times" ("end_date");

CREATE INDEX ON "event_check_in_steps" ("created_at");

CREATE INDEX ON "event_check_in_steps" ("id");

CREATE INDEX ON "event_check_in_steps" ("event_date_time_id");

CREATE INDEX ON "event_date_details" ("event_date_time_id");

CREATE INDEX ON "event_date_details" ("start_time");

CREATE INDEX ON "event_date_details" ("end_time");

CREATE INDEX ON "event_date_details" ("created_at");

CREATE INDEX ON "event_date_publishes" ("event_date_time_id");

CREATE INDEX ON "event_date_publishes" ("event_public");

CREATE INDEX ON "event_date_publishes" ("event_going_public_time");

CREATE INDEX ON "event_date_publishes" ("created_at");

CREATE INDEX ON "event_date_private_audiences" ("id");

CREATE INDEX ON "event_date_private_audiences" ("event_date_time_id");

CREATE INDEX ON "event_date_private_audiences" ("created_at");

CREATE INDEX ON "event_date_tickets" ("id");

CREATE INDEX ON "event_date_tickets" ("event_date_time_id");

CREATE INDEX ON "event_date_tickets" ("start_date");

CREATE INDEX ON "event_date_tickets" ("end_date");

CREATE INDEX ON "event_date_tickets" ("name");

CREATE INDEX ON "event_date_tickets" ("price");

CREATE INDEX ON "event_date_tickets" ("absorb_fees");

CREATE INDEX ON "event_date_tickets" ("capacity");

CREATE INDEX ON "event_date_tickets" ("type");

CREATE INDEX ON "event_date_tickets" ("ticket_type");

CREATE INDEX ON "event_date_tickets" ("created_at");

CREATE INDEX ON "event_date_locations" ("city");

CREATE INDEX ON "event_date_locations" ("state");

CREATE INDEX ON "event_date_locations" ("country");

CREATE INDEX ON "event_date_locations" ("event_date_time_id");

CREATE INDEX ON "event_date_locations" ("created_at");

CREATE INDEX ON "locations" ("option_id");

CREATE INDEX ON "locations" ("city");

CREATE INDEX ON "locations" ("state");

CREATE INDEX ON "locations" ("country");

CREATE INDEX ON "locations" ("created_at");

CREATE INDEX ON "locations" ("updated_at");

CREATE INDEX ON "shortlets" ("space_type");

CREATE INDEX ON "shortlets" ("type_of_shortlet");

CREATE INDEX ON "shortlets" ("option_id");

CREATE INDEX ON "shortlets" ("created_at");

CREATE INDEX ON "shortlets" ("updated_at");

CREATE INDEX ON "option_date_times" ("id");

CREATE INDEX ON "option_date_times" ("available");

CREATE INDEX ON "option_date_times" ("price");

CREATE INDEX ON "option_date_times" ("option_id");

CREATE INDEX ON "option_date_times" ("created_at");

CREATE INDEX ON "option_date_times" ("updated_at");

CREATE INDEX ON "space_areas" ("id");

CREATE INDEX ON "space_areas" ("space_type");

CREATE INDEX ON "space_areas" ("option_id");

CREATE INDEX ON "space_areas" ("created_at");

CREATE INDEX ON "space_areas" ("updated_at");

CREATE INDEX ON "options_prices" ("option_id");

CREATE INDEX ON "options_prices" ("created_at");

CREATE INDEX ON "wifi_details" ("option_id");

CREATE INDEX ON "wifi_details" ("created_at");

CREATE INDEX ON "wifi_details" ("updated_at");

CREATE INDEX ON "amenities" ("id");

CREATE INDEX ON "amenities" ("tag");

CREATE INDEX ON "amenities" ("am_type");

CREATE INDEX ON "amenities" ("size_option");

CREATE INDEX ON "amenities" ("has_am");

CREATE INDEX ON "amenities" ("option_id");

CREATE INDEX ON "amenities" ("created_at");

CREATE INDEX ON "amenities" ("updated_at");

CREATE INDEX ON "check_in_steps" ("created_at");

CREATE INDEX ON "check_in_steps" ("id");

CREATE INDEX ON "check_in_steps" ("option_id");

CREATE INDEX ON "option_availability_settings" ("created_at");

CREATE INDEX ON "option_availability_settings" ("option_id");

CREATE INDEX ON "cancel_policies" ("created_at");

CREATE INDEX ON "cancel_policies" ("option_id");

CREATE INDEX ON "option_trip_lengths" ("created_at");

CREATE INDEX ON "option_trip_lengths" ("option_id");

CREATE INDEX ON "check_in_out_details" ("created_at");

CREATE INDEX ON "check_in_out_details" ("option_id");

CREATE INDEX ON "option_co_hosts" ("created_at");

CREATE INDEX ON "option_co_hosts" ("option_id");

CREATE INDEX ON "option_co_hosts" ("id");

CREATE INDEX ON "options_extra_infos" ("id");

CREATE INDEX ON "options_extra_infos" ("option_id");

CREATE INDEX ON "options_extra_infos" ("type");

CREATE INDEX ON "options_extra_infos" ("created_at");

CREATE INDEX ON "options_extra_infos" ("updated_at");

CREATE INDEX ON "things_to_notes" ("id");

CREATE INDEX ON "things_to_notes" ("option_id");

CREATE INDEX ON "things_to_notes" ("type");

CREATE INDEX ON "things_to_notes" ("tag");

CREATE INDEX ON "things_to_notes" ("created_at");

CREATE INDEX ON "things_to_notes" ("updated_at");

CREATE INDEX ON "option_invitations" ("id");

CREATE INDEX ON "option_invitations" ("option_id");

CREATE INDEX ON "option_invitations" ("created_at");

CREATE INDEX ON "option_rules" ("id");

CREATE INDEX ON "option_rules" ("option_id");

CREATE INDEX ON "option_rules" ("type");

CREATE INDEX ON "option_rules" ("tag");

CREATE INDEX ON "option_rules" ("created_at");

CREATE INDEX ON "option_rules" ("updated_at");

CREATE INDEX ON "option_book_methods" ("option_id");

CREATE INDEX ON "option_book_methods" ("instant_book");

CREATE INDEX ON "option_book_methods" ("created_at");

CREATE INDEX ON "option_book_methods" ("updated_at");

CREATE INDEX ON "book_requirements" ("option_id");

CREATE INDEX ON "book_requirements" ("profile_photo");

CREATE INDEX ON "book_requirements" ("created_at");

CREATE INDEX ON "book_requirements" ("updated_at");

CREATE INDEX ON "option_messages" ("option_id");

CREATE INDEX ON "option_messages" ("id");

CREATE INDEX ON "option_messages" ("user_id");

CREATE INDEX ON "option_messages" ("seen");

CREATE INDEX ON "option_messages" ("created_at");

CREATE INDEX ON "option_messages" ("updated_at");

CREATE INDEX ON "feedbacks" ("id");

CREATE INDEX ON "feedbacks" ("user_id");

CREATE INDEX ON "feedbacks" ("created_at");

CREATE INDEX ON "feedbacks" ("updated_at");

CREATE INDEX ON "option_reference_infos" ("option_charge_id");

CREATE INDEX ON "option_reference_infos" ("created_at");

CREATE INDEX ON "event_reference_infos" ("event_charge_id");

CREATE INDEX ON "event_reference_infos" ("created_at");

CREATE INDEX ON "messages" ("id");

CREATE INDEX ON "messages" ("sender_id");

CREATE INDEX ON "messages" ("receiver_id");

CREATE INDEX ON "messages" ("type");

CREATE INDEX ON "messages" ("created_at");

CREATE INDEX ON "request_notifies" ("m_id");

CREATE INDEX ON "request_notifies" ("created_at");

CREATE INDEX ON "single_rooms" ("id");

CREATE INDEX ON "single_rooms" ("user_one");

CREATE INDEX ON "single_rooms" ("user_two");

CREATE INDEX ON "notifications" ("id");

ALTER TABLE "user_apn_details" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "users_profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "users_options_reviews" ADD FOREIGN KEY ("option_user_id") REFERENCES "options_infos" ("option_user_id");

ALTER TABLE "users_options_reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "identity" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "em_contacts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "mailing_addresses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "account_numbers" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "payments_gate_pays" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "cards" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "users_locations" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("id") REFERENCES "users" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "option_questions" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "wishlists" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wishlists_items" ADD FOREIGN KEY ("wishlist_id") REFERENCES "wishlists" ("id");

ALTER TABLE "wishlists_items" ADD FOREIGN KEY ("option_user_id") REFERENCES "options_infos" ("option_user_id");

ALTER TABLE "vids" ADD FOREIGN KEY ("option_user_id") REFERENCES "options_infos" ("option_user_id");

ALTER TABLE "vids" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "options_infos" ADD FOREIGN KEY ("host_id") REFERENCES "users" ("id");

ALTER TABLE "options_infos" ADD FOREIGN KEY ("primary_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "options_infos_category" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "events_infos_category" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "options_infos_status" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "complete_option_info" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "options_info_details" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_add_charges" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_discounts" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "options_info_photos" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "options_photo_captions" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "event_infos" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "event_date_times" ADD FOREIGN KEY ("event_info_id") REFERENCES "event_infos" ("option_id");

ALTER TABLE "event_check_in_steps" ADD FOREIGN KEY ("event_date_time_id") REFERENCES "event_date_times" ("id");

ALTER TABLE "event_date_details" ADD FOREIGN KEY ("event_date_time_id") REFERENCES "event_date_times" ("id");

ALTER TABLE "event_date_publishes" ADD FOREIGN KEY ("event_date_time_id") REFERENCES "event_date_times" ("id");

ALTER TABLE "event_date_private_audiences" ADD FOREIGN KEY ("event_date_time_id") REFERENCES "event_date_publishes" ("event_date_time_id");

ALTER TABLE "event_date_tickets" ADD FOREIGN KEY ("event_date_time_id") REFERENCES "event_date_times" ("id");

ALTER TABLE "event_date_locations" ADD FOREIGN KEY ("event_date_time_id") REFERENCES "event_date_times" ("id");

ALTER TABLE "locations" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "shortlets" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_date_times" ADD FOREIGN KEY ("option_id") REFERENCES "shortlets" ("option_id");

ALTER TABLE "space_areas" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "options_prices" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "wifi_details" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "amenities" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "check_in_steps" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_availability_settings" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "cancel_policies" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_trip_lengths" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "check_in_out_details" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_co_hosts" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "options_extra_infos" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "things_to_notes" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_invitations" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_rules" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_book_methods" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "book_requirements" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_messages" ADD FOREIGN KEY ("option_id") REFERENCES "options_infos" ("id");

ALTER TABLE "option_messages" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "feedbacks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "report_options" ADD FOREIGN KEY ("option_user_id") REFERENCES "options_infos" ("option_user_id");

ALTER TABLE "report_options" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "charge_references" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "cancellations" ADD FOREIGN KEY ("cancel_user_id") REFERENCES "users" ("user_id");

ALTER TABLE "refunds" ADD FOREIGN KEY ("charge_id") REFERENCES "main_refunds" ("charge_id");

ALTER TABLE "refunds" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "refund_payouts" ADD FOREIGN KEY ("charge_id") REFERENCES "main_refunds" ("charge_id");

ALTER TABLE "refund_payouts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "payouts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "charge_option_references" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "charge_option_references" ADD FOREIGN KEY ("option_user_id") REFERENCES "options_infos" ("option_user_id");

ALTER TABLE "charge_event_references" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "charge_event_references" ADD FOREIGN KEY ("option_user_id") REFERENCES "options_infos" ("option_user_id");

ALTER TABLE "charge_date_references" ADD FOREIGN KEY ("charge_event_id") REFERENCES "charge_event_references" ("id");

ALTER TABLE "charge_date_references" ADD FOREIGN KEY ("event_date_id") REFERENCES "event_date_times" ("id");

ALTER TABLE "charge_ticket_references" ADD FOREIGN KEY ("charge_date_id") REFERENCES "charge_date_references" ("id");

ALTER TABLE "charge_ticket_references" ADD FOREIGN KEY ("ticket_id") REFERENCES "event_date_tickets" ("id");

ALTER TABLE "complete_charge_reviews" ADD FOREIGN KEY ("charge_review_id") REFERENCES "charge_reviews" ("charge_id");

ALTER TABLE "scanned_charges" ADD FOREIGN KEY ("scanned_by") REFERENCES "users" ("user_id");

ALTER TABLE "option_reference_infos" ADD FOREIGN KEY ("option_charge_id") REFERENCES "charge_option_references" ("id");

ALTER TABLE "event_reference_infos" ADD FOREIGN KEY ("event_charge_id") REFERENCES "charge_event_references" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("sender_id") REFERENCES "users" ("user_id");

ALTER TABLE "messages" ADD FOREIGN KEY ("receiver_id") REFERENCES "users" ("user_id");

ALTER TABLE "request_notifies" ADD FOREIGN KEY ("m_id") REFERENCES "messages" ("id");

ALTER TABLE "single_rooms" ADD FOREIGN KEY ("user_one") REFERENCES "users" ("user_id");

ALTER TABLE "single_rooms" ADD FOREIGN KEY ("user_two") REFERENCES "users" ("user_id");

ALTER TABLE "notifications" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");