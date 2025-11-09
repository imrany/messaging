-- =============================================
-- DOWN Migration - Drop all tables
-- =============================================

-- Drop realtime
ALTER PUBLICATION supabase_realtime DROP TABLE public.alerts;
ALTER PUBLICATION supabase_realtime DROP TABLE public.sensor_readings;

-- Drop indexes
DROP INDEX IF EXISTS idx_alerts_resolved;
DROP INDEX IF EXISTS idx_alerts_hub_id;
DROP INDEX IF EXISTS idx_sensor_readings_recorded_at;
DROP INDEX IF EXISTS idx_sensor_readings_hub_id;

-- Drop tables in reverse order (respecting foreign keys)
DROP TABLE IF EXISTS public.alerts;
DROP TABLE IF EXISTS public.sensor_readings;
DROP TABLE IF EXISTS public.notification_preferences;
DROP TABLE IF EXISTS public.market_listings;
DROP TABLE IF EXISTS public.courses;
DROP TABLE IF EXISTS public.hubs;
DROP TABLE IF EXISTS public.profiles;

-- Drop enum
DROP TYPE IF EXISTS public.user_role;
