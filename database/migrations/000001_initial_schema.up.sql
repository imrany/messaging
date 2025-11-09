-- =============================================
-- UP Migration - Create all tables
-- =============================================

-- Create enum for user roles
CREATE TYPE user_role AS ENUM ('farmer', 'buyer', 'admin');

-- Create profiles table
CREATE TABLE profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  full_name TEXT NOT NULL,
  phone TEXT UNIQUE,
  password TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  role user_role NOT NULL,
  location TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create hubs table
CREATE TABLE hubs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  location TEXT NOT NULL,
  manager_id UUID REFERENCES profiles(id),
  description TEXT,
  contact_phone TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create courses table
CREATE TABLE courses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  description TEXT,
  content TEXT,
  image_url TEXT,
  duration TEXT,
  level TEXT DEFAULT 'beginner',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create market listings table
CREATE TABLE market_listings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  farmer_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
  product_name TEXT NOT NULL,
  description TEXT,
  quantity DECIMAL NOT NULL,
  unit TEXT NOT NULL,
  price_per_unit DECIMAL NOT NULL,
  available BOOLEAN DEFAULT true,
  image_url TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create notification preferences table
CREATE TABLE notification_preferences (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
  sms_enabled BOOLEAN DEFAULT false,
  whatsapp_enabled BOOLEAN DEFAULT false,
  email_enabled BOOLEAN DEFAULT false,
  phone_number TEXT,
  email TEXT REFERENCES profiles(email) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create sensor readings table
CREATE TABLE sensor_readings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hub_id UUID REFERENCES hubs(id) ON DELETE CASCADE NOT NULL,
  temperature DECIMAL(5,2) NOT NULL,
  humidity DECIMAL(5,2) NOT NULL,
  recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create alerts table
CREATE TABLE alerts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hub_id UUID REFERENCES hubs(id) ON DELETE CASCADE NOT NULL,
  alert_type TEXT NOT NULL CHECK (alert_type IN ('temperature', 'humidity', 'both')),
  message TEXT NOT NULL,
  temperature DECIMAL(5,2),
  humidity DECIMAL(5,2),
  resolved BOOLEAN DEFAULT false,
  resolved_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create indexes
CREATE INDEX idx_sensor_readings_hub_id ON sensor_readings(hub_id);
CREATE INDEX idx_sensor_readings_recorded_at ON sensor_readings(recorded_at DESC);
CREATE INDEX idx_alerts_hub_id ON alerts(hub_id);
CREATE INDEX idx_alerts_resolved ON alerts(resolved);

-- Enable RLS on all tables
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE hubs ENABLE ROW LEVEL SECURITY;
ALTER TABLE courses ENABLE ROW LEVEL SECURITY;
ALTER TABLE market_listings ENABLE ROW LEVEL SECURITY;
ALTER TABLE notification_preferences ENABLE ROW LEVEL SECURITY;
ALTER TABLE sensor_readings ENABLE ROW LEVEL SECURITY;
ALTER TABLE alerts ENABLE ROW LEVEL SECURITY;
