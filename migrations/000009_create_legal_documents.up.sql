ALTER TABLE system_logs
DROP COLUMN category;

ALTER TABLE system_logs
ADD COLUMN category ENUM('SECURITY', 'SYSTEM', 'LEGAL', 'CONSENT') NOT NULL;

-- 2. Create the Legal Documents table (The "Source of Truth")
CREATE TABLE IF NOT EXISTS legal_documents (
    id INT AUTO_INCREMENT PRIMARY KEY,
    doc_type ENUM('PRIVACY_POLICY', 'TERMS_OF_SERVICE') NOT NULL, -- e.g., 'PRIVACY_POLICY'
    version VARCHAR(20) NOT NULL,  -- e.g., '2026.03.10'
    file_url TEXT NOT NULL,       -- Link to the file in Azure Storage
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doc_type, version)      -- Prevents duplicate version numbers
);

CREATE TABLE IF NOT EXISTS user_consents (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_email VARCHAR(255) NOT NULL,
    document_id INT NOT NULL,
    accepted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45), -- To log the user's IP address at the time of consent
    FOREIGN KEY (document_id) REFERENCES legal_documents(id),
    UNIQUE KEY unique_user_consent (user_email, document_id)
)