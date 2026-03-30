-- AI Gateway 数据库表结构

-- 租户表
CREATE TABLE `tenants` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '租户名称',
  `status` int DEFAULT 1 COMMENT '状态：1启用 0禁用',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- API Key表
CREATE TABLE `api_keys` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tenant_id` bigint unsigned NOT NULL COMMENT '租户ID',
  `key_hash` varchar(64) NOT NULL COMMENT 'API密钥哈希',
  `key_prefix` varchar(20) DEFAULT NULL COMMENT 'API密钥前缀',
  `name` varchar(100) DEFAULT NULL COMMENT '密钥名称',
  `provider_id` bigint unsigned DEFAULT NULL COMMENT 'Provider ID',
  `status` int DEFAULT 1 COMMENT '状态：1启用 0禁用',
  `rate_limit_per_sec` int DEFAULT 10 COMMENT '每秒请求限制，0为不限',
  `rate_limit` int DEFAULT 60 COMMENT '每分钟请求限制',
  `monthly_token_limit` int DEFAULT 0 COMMENT '每月Token限额，0为不限',
  `total_token_quota` int DEFAULT 0 COMMENT '总Token额度，0为不限',
  `tokens_consumed` int DEFAULT 0 COMMENT '已消耗Token数',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_api_keys_key_hash` (`key_hash`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_provider_id` (`provider_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API密钥表';

-- 请求日志表
CREATE TABLE `request_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tenant_id` bigint unsigned NOT NULL COMMENT '租户ID',
  `api_key_id` bigint unsigned DEFAULT NULL COMMENT 'API密钥ID',
  `provider_id` bigint unsigned DEFAULT NULL COMMENT 'Provider ID',
  `model` varchar(50) DEFAULT NULL COMMENT '模型名称',
  `tokens_used` int DEFAULT 0 COMMENT '总Token数',
  `prompt_tokens` int DEFAULT 0 COMMENT '输入Token数',
  `comp_tokens` int DEFAULT 0 COMMENT '输出Token数',
  `status` int DEFAULT NULL COMMENT '状态码',
  `latency` int DEFAULT NULL COMMENT '延迟(ms)',
  `ip` varchar(50) DEFAULT NULL COMMENT '请求IP',
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_api_key_id` (`api_key_id`),
  KEY `idx_provider_id` (`provider_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='请求日志表';

-- Token使用统计表
CREATE TABLE `token_usages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tenant_id` bigint unsigned NOT NULL COMMENT '租户ID',
  `api_key_id` bigint unsigned DEFAULT NULL COMMENT 'API密钥ID',
  `tokens` int DEFAULT NULL COMMENT 'Token数量',
  `date` datetime DEFAULT NULL COMMENT '日期',
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_apikey_date` (`tenant_id`, `api_key_id`, `date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Token使用统计表';

-- Provider表
CREATE TABLE `providers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT 'Provider名称',
  `type` varchar(20) NOT NULL COMMENT 'Provider类型',
  `base_url` varchar(255) DEFAULT NULL COMMENT 'Base URL',
  `api_key_enc` varchar(500) DEFAULT NULL COMMENT '加密的API密钥',
  `models` text COMMENT '支持的模型列表',
  `status` int DEFAULT 1 COMMENT '状态：1启用 0禁用',
  `priority` int DEFAULT 0 COMMENT '优先级',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Provider表';
