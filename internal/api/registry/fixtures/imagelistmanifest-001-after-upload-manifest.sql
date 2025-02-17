INSERT INTO accounts (name, auth_tenant_id) VALUES ('test1', 'test1authtenant');

INSERT INTO blob_mounts (blob_id, repo_id) VALUES (1, 1);
INSERT INTO blob_mounts (blob_id, repo_id) VALUES (2, 1);
INSERT INTO blob_mounts (blob_id, repo_id) VALUES (3, 1);
INSERT INTO blob_mounts (blob_id, repo_id) VALUES (4, 1);

INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, media_type) VALUES (1, 'test1', 'sha256:442f91fa9998460f28e8ff7023e5ddca679f7d2b51dc5498e8aba249678cc7f8', 1048919, '6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b', 1, 1, 'application/vnd.docker.image.rootfs.diff.tar.gzip');
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, media_type) VALUES (2, 'test1', 'sha256:a0a84c915810634c0d4522dca789fa95a7ad5b843860ead04d2e13ec949d8a2f', 1257, 'd4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35', 1, 1, 'application/vnd.docker.container.image.v1+json');
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, media_type) VALUES (3, 'test1', 'sha256:3ae14a50df760250f0e97faf429cc4541c832ed0de61ad5b6ac25d1d695d1a6e', 1048919, '4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce', 2, 2, 'application/vnd.docker.image.rootfs.diff.tar.gzip');
INSERT INTO blobs (id, account_name, digest, size_bytes, storage_id, pushed_at, validated_at, media_type) VALUES (4, 'test1', 'sha256:958a896c0ef1220adf55b23d10e8e960a658b451ace48597f44db27a4a899304', 1257, '4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a', 2, 2, 'application/vnd.docker.container.image.v1+json');

INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:4c4f2bca300e74786a04590aa15cfcbfa1f3ec64c15fad0a0df8a6674dcbf34b', 3);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:4c4f2bca300e74786a04590aa15cfcbfa1f3ec64c15fad0a0df8a6674dcbf34b', 4);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:e3c1e46560a7ce30e3d107791e1f60a588eda9554564a5d17aa365e53dd6ae58', 1);
INSERT INTO manifest_blob_refs (repo_id, digest, blob_id) VALUES (1, 'sha256:e3c1e46560a7ce30e3d107791e1f60a588eda9554564a5d17aa365e53dd6ae58', 2);

INSERT INTO manifest_contents (repo_id, digest, content) VALUES (1, 'sha256:4c4f2bca300e74786a04590aa15cfcbfa1f3ec64c15fad0a0df8a6674dcbf34b', '{"config":{"digest":"sha256:958a896c0ef1220adf55b23d10e8e960a658b451ace48597f44db27a4a899304","mediaType":"application/vnd.docker.container.image.v1+json","size":1257},"layers":[{"digest":"sha256:3ae14a50df760250f0e97faf429cc4541c832ed0de61ad5b6ac25d1d695d1a6e","mediaType":"application/vnd.docker.image.rootfs.diff.tar.gzip","size":1048919}],"mediaType":"application/vnd.docker.distribution.manifest.v2+json","schemaVersion":2}');
INSERT INTO manifest_contents (repo_id, digest, content) VALUES (1, 'sha256:dc8b0fc112e08d16a5d1b608ab928aea0a6f5484b8c17ee06afa825a75eadc44', '{"manifests":[{"digest":"sha256:e3c1e46560a7ce30e3d107791e1f60a588eda9554564a5d17aa365e53dd6ae58","mediaType":"application/vnd.docker.distribution.manifest.v2+json","platform":{"architecture":"amd64","os":"linux"},"size":428},{"digest":"sha256:4c4f2bca300e74786a04590aa15cfcbfa1f3ec64c15fad0a0df8a6674dcbf34b","mediaType":"application/vnd.docker.distribution.manifest.v2+json","platform":{"architecture":"arm","os":"linux"},"size":428}],"mediaType":"application/vnd.docker.distribution.manifest.list.v2+json","schemaVersion":2}');
INSERT INTO manifest_contents (repo_id, digest, content) VALUES (1, 'sha256:e3c1e46560a7ce30e3d107791e1f60a588eda9554564a5d17aa365e53dd6ae58', '{"config":{"digest":"sha256:a0a84c915810634c0d4522dca789fa95a7ad5b843860ead04d2e13ec949d8a2f","mediaType":"application/vnd.docker.container.image.v1+json","size":1257},"layers":[{"digest":"sha256:442f91fa9998460f28e8ff7023e5ddca679f7d2b51dc5498e8aba249678cc7f8","mediaType":"application/vnd.docker.image.rootfs.diff.tar.gzip","size":1048919}],"mediaType":"application/vnd.docker.distribution.manifest.v2+json","schemaVersion":2}');

INSERT INTO manifest_manifest_refs (repo_id, parent_digest, child_digest) VALUES (1, 'sha256:dc8b0fc112e08d16a5d1b608ab928aea0a6f5484b8c17ee06afa825a75eadc44', 'sha256:4c4f2bca300e74786a04590aa15cfcbfa1f3ec64c15fad0a0df8a6674dcbf34b');
INSERT INTO manifest_manifest_refs (repo_id, parent_digest, child_digest) VALUES (1, 'sha256:dc8b0fc112e08d16a5d1b608ab928aea0a6f5484b8c17ee06afa825a75eadc44', 'sha256:e3c1e46560a7ce30e3d107791e1f60a588eda9554564a5d17aa365e53dd6ae58');

INSERT INTO manifests (repo_id, digest, media_type, size_bytes, pushed_at, validated_at) VALUES (1, 'sha256:4c4f2bca300e74786a04590aa15cfcbfa1f3ec64c15fad0a0df8a6674dcbf34b', 'application/vnd.docker.distribution.manifest.v2+json', 1050604, 2, 2);
INSERT INTO manifests (repo_id, digest, media_type, size_bytes, pushed_at, validated_at) VALUES (1, 'sha256:dc8b0fc112e08d16a5d1b608ab928aea0a6f5484b8c17ee06afa825a75eadc44', 'application/vnd.docker.distribution.manifest.list.v2+json', 2101735, 3, 3);
INSERT INTO manifests (repo_id, digest, media_type, size_bytes, pushed_at, validated_at) VALUES (1, 'sha256:e3c1e46560a7ce30e3d107791e1f60a588eda9554564a5d17aa365e53dd6ae58', 'application/vnd.docker.distribution.manifest.v2+json', 1050604, 1, 1);

INSERT INTO quotas (auth_tenant_id, manifests) VALUES ('test1authtenant', 100);

INSERT INTO repos (id, account_name, name) VALUES (1, 'test1', 'foo');

INSERT INTO tags (repo_id, name, digest, pushed_at) VALUES (1, 'first', 'sha256:e3c1e46560a7ce30e3d107791e1f60a588eda9554564a5d17aa365e53dd6ae58', 1);
INSERT INTO tags (repo_id, name, digest, pushed_at) VALUES (1, 'list', 'sha256:dc8b0fc112e08d16a5d1b608ab928aea0a6f5484b8c17ee06afa825a75eadc44', 3);
INSERT INTO tags (repo_id, name, digest, pushed_at) VALUES (1, 'second', 'sha256:4c4f2bca300e74786a04590aa15cfcbfa1f3ec64c15fad0a0df8a6674dcbf34b', 2);

INSERT INTO vuln_info (repo_id, digest, status, message, next_check_at) VALUES (1, 'sha256:4c4f2bca300e74786a04590aa15cfcbfa1f3ec64c15fad0a0df8a6674dcbf34b', 'Pending', '', 2);
INSERT INTO vuln_info (repo_id, digest, status, message, next_check_at) VALUES (1, 'sha256:dc8b0fc112e08d16a5d1b608ab928aea0a6f5484b8c17ee06afa825a75eadc44', 'Pending', '', 3);
INSERT INTO vuln_info (repo_id, digest, status, message, next_check_at) VALUES (1, 'sha256:e3c1e46560a7ce30e3d107791e1f60a588eda9554564a5d17aa365e53dd6ae58', 'Pending', '', 1);
