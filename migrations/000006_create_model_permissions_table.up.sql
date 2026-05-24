-- Polymorphic pivot: users/admins -> permissions
CREATE TABLE IF NOT EXISTS `model_permissions` (
  `model_id` bigint unsigned NOT NULL,
  `model_type` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`model_id`, `model_type`, `permission_id`),
  KEY `idx_model_permissions_model_id` (`model_id`),
  KEY `idx_model_permissions_model_type` (`model_type`),
  KEY `idx_model_permissions_permission_id` (`permission_id`),
  CONSTRAINT `fk_model_permissions_permission_id`
    FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
