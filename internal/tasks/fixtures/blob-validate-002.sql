INSERT INTO accounts (name, auth_tenant_id, upstream_peer_hostname, required_labels, blobs_sweeped_at, storage_sweeped_at, metadata_json) VALUES ('test1', 'test1authtenant', '', '', NULL, NULL, '');

INSERT INTO blob_mounts (blob_id, repo_id, marked_for_deletion_at) VALUES (1, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, marked_for_deletion_at) VALUES (2, 1, NULL);
INSERT INTO blob_mounts (blob_id, repo_id, marked_for_deletion_at) VALUES (3, 1, NULL);

INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, marked_for_deletion_at) VALUES (1, 'test1', 'sha256:3d309a098968afe810ee167e0c5b205ef3610829a6d34e0f0ba4ca66756c6f5e', 1048576, '3d309a098968afe810ee167e0c5b205ef3610829a6d34e0f0ba4ca66756c6f5e', 3601, 1386001, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, marked_for_deletion_at) VALUES (2, 'test1', 'sha256:a718f4a112454b50c8ecd2b0a5b00eb32ee90699593625139cd3fabc97dcce8d', 1048576, 'a718f4a112454b50c8ecd2b0a5b00eb32ee90699593625139cd3fabc97dcce8d', 3602, 1386002, '', NULL);
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, validation_error_message, marked_for_deletion_at) VALUES (3, 'test1', 'sha256:5aa3b2bb09eef96742be217d3ef99cbde12ee2817a1e6d5d15de5922354aba40', 1048576, 'd87b0830e55e19fd0825bebaa110ebade6de7d1fcf2ddf0fbd5762e1f809370e', 3603, 1386003, 'expected digest sha256:5aa3b2bb09eef96742be217d3ef99cbde12ee2817a1e6d5d15de5922354aba40, but got sha256:d87b0830e55e19fd0825bebaa110ebade6de7d1fcf2ddf0fbd5762e1f809370e', NULL);

INSERT INTO repos (id, account_name, name, blob_mounts_sweeped_at, manifests_synced_at) VALUES (1, 'test1', 'foo', NULL, NULL);
