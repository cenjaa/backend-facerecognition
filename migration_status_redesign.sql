-- ============================================
-- Migration: Schema Redesign to match new ERD
-- Run this on the LIVE VPS PostgreSQL database
-- ============================================

-- Step 1: Rename ms_user.name → nama (if column still called 'name')
ALTER TABLE ms_user RENAME COLUMN name TO nama;

-- Step 2: Add new columns to ms_user
ALTER TABLE ms_user ADD COLUMN IF NOT EXISTS created_by VARCHAR(255) DEFAULT '';
ALTER TABLE ms_user ADD COLUMN IF NOT EXISTS dataset_path VARCHAR(500) DEFAULT '';
ALTER TABLE ms_user ADD COLUMN IF NOT EXISTS akumulasi_kpi INTEGER NOT NULL DEFAULT 0;

-- Step 3: Migrate akumulasi data from tbl_akumulasi_kpi_harian into ms_user
UPDATE ms_user u SET akumulasi_kpi = COALESCE((
    SELECT SUM(a.akumulasi_kpi) FROM tbl_akumulasi_kpi_harian a WHERE a.id_user = u.id_user
), 0);

-- Step 4: Drop the old tbl_akumulasi_kpi_harian table
DROP TABLE IF EXISTS tbl_akumulasi_kpi_harian;

-- Step 5: Add username/password columns to ms_admin
ALTER TABLE ms_admin ADD COLUMN IF NOT EXISTS username VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE ms_admin ADD COLUMN IF NOT EXISTS password VARCHAR(255) NOT NULL DEFAULT '';

-- Step 6: Rename ms_status.description → deskripsi (if column still called 'description')
ALTER TABLE ms_status RENAME COLUMN description TO deskripsi;

-- Step 7: Add story_jira column to tbl_kpi_harian
ALTER TABLE tbl_kpi_harian ADD COLUMN IF NOT EXISTS story_jira VARCHAR(255) DEFAULT '';

-- Step 8: Rename tbl_kpi_harian.nama → nama_tiket
ALTER TABLE tbl_kpi_harian RENAME COLUMN nama TO nama_tiket;

-- Step 9: Fix ms_status seed data with new column name
INSERT INTO ms_status (id_status, deskripsi) VALUES (0, 'Menunggu Kehadiran') ON CONFLICT (id_status) DO NOTHING;
INSERT INTO ms_status (id_status, deskripsi) VALUES (3, 'Telat Hadir') ON CONFLICT (id_status) DO UPDATE SET deskripsi = 'Telat Hadir';
INSERT INTO ms_status (id_status, deskripsi) VALUES (4, 'Keluar') ON CONFLICT (id_status) DO NOTHING;
INSERT INTO ms_status (id_status, deskripsi) VALUES (5, 'Cuti') ON CONFLICT (id_status) DO NOTHING;

-- Verify
SELECT * FROM ms_status ORDER BY id_status;
SELECT column_name FROM information_schema.columns WHERE table_name = 'ms_user' ORDER BY ordinal_position;
SELECT column_name FROM information_schema.columns WHERE table_name = 'ms_admin' ORDER BY ordinal_position;
SELECT column_name FROM information_schema.columns WHERE table_name = 'tbl_kpi_harian' ORDER BY ordinal_position;
