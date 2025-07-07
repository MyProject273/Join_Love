-- Refactored schema using UUIDs for all primary and foreign keys

CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- Needed for gen_random_uuid()

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_name" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" varchar UNIQUE,
  "password_hash" text NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user',
  "full_name" varchar,
  "gender" varchar,
  "birthdate" date,
  "avatar_url" text,
  "bio" text,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_verified" boolean NOT NULL DEFAULT false,
  "last_login" timestamp,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "user_profiles" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID UNIQUE NOT NULL,
  "job_title" varchar,
  "education" varchar,
  "interests" text,
  "relationship_goal" varchar,
  "religion" varchar,
  "smoking" varchar,
  "drinking" varchar,
  "height_cm" int,
  "weight_kg" int,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "user_preferences" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID UNIQUE NOT NULL,
  "preferred_gender" varchar,
  "min_age" int,
  "max_age" int,
  "max_distance_km" int,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "photos" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "image_url" text NOT NULL,
  "is_main" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "user_locations" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID UNIQUE NOT NULL,
  "latitude" decimal,
  "longitude" decimal,
  "city" varchar,
  "country" varchar,
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "created_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "likes" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "from_user_id" UUID NOT NULL,
  "to_user_id" UUID NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "matches" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user1_id" UUID NOT NULL,
  "user2_id" UUID NOT NULL,
  "matched_at" timestamp NOT NULL DEFAULT now(),
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "messages" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "match_id" UUID NOT NULL,
  "sender_id" UUID NOT NULL,
  "content" text NOT NULL,
  "sent_at" timestamp NOT NULL DEFAULT now(),
  "is_read" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "message_attachments" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "message_id" UUID NOT NULL,
  "file_url" text NOT NULL,
  "file_type" varchar NOT NULL,
  "uploaded_at" timestamp NOT NULL DEFAULT now(),
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "blocks" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "blocker_id" UUID NOT NULL,
  "blocked_id" UUID NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "reports" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "reporter_id" UUID NOT NULL,
  "reported_id" UUID NOT NULL,
  "reason" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "privacy_settings" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),  
  "user_id" UUID UNIQUE NOT NULL,
  "show_age" boolean NOT NULL DEFAULT true,
  "show_distance" boolean NOT NULL DEFAULT true,
  "allow_messages_from_non_matches" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "notifications" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "type" varchar NOT NULL,
  "message" text NOT NULL,
  "is_read" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "user_activities" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "action" varchar NOT NULL,
  "metadata" text,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "subscriptions" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "plan_name" varchar NOT NULL,
  "price" decimal NOT NULL,
  "started_at" timestamp NOT NULL,
  "expires_at" timestamp NOT NULL,
  "is_active" boolean NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

CREATE TABLE "payments" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "subscription_id" UUID NOT NULL,
  "amount" decimal NOT NULL,
  "payment_method" varchar NOT NULL,
  "payment_status" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  "deleted_at" timestamp,
  "created_by" UUID
);

-- Foreign keys
ALTER TABLE "user_profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_preferences" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "photos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_locations" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "likes" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");
ALTER TABLE "likes" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");
ALTER TABLE "matches" ADD FOREIGN KEY ("user1_id") REFERENCES "users" ("id");
ALTER TABLE "matches" ADD FOREIGN KEY ("user2_id") REFERENCES "users" ("id");
ALTER TABLE "messages" ADD FOREIGN KEY ("match_id") REFERENCES "matches" ("id");
ALTER TABLE "messages" ADD FOREIGN KEY ("sender_id") REFERENCES "users" ("id");
ALTER TABLE "message_attachments" ADD FOREIGN KEY ("message_id") REFERENCES "messages" ("id");
ALTER TABLE "blocks" ADD FOREIGN KEY ("blocker_id") REFERENCES "users" ("id");
ALTER TABLE "blocks" ADD FOREIGN KEY ("blocked_id") REFERENCES "users" ("id");
ALTER TABLE "reports" ADD FOREIGN KEY ("reporter_id") REFERENCES "users" ("id");
ALTER TABLE "reports" ADD FOREIGN KEY ("reported_id") REFERENCES "users" ("id");
ALTER TABLE "privacy_settings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "notifications" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_activities" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "subscriptions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "payments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "payments" ADD FOREIGN KEY ("subscription_id") REFERENCES "subscriptions" ("id");
