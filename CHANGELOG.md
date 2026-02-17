# Changelog

## 0.1.0-alpha.71 (2026-02-17)

Full Changelog: [v0.1.0-alpha.70...v0.1.0-alpha.71](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.70...v0.1.0-alpha.71)

### ⚠ BREAKING CHANGES

* add support for passing files as parameters

### Features

* add exponential backoff retry to PullOutput ([d218b14](https://github.com/stainless-api/stainless-api-cli/commit/d218b14d95fef2b12a65eabe6d7fa0df37e834f4))
* add github version check ([3596cbb](https://github.com/stainless-api/stainless-api-cli/commit/3596cbba986e3019e9fc1f9b555728a5c737fe27))
* add readme documentation for passing files as arguments ([3d5c3eb](https://github.com/stainless-api/stainless-api-cli/commit/3d5c3eb26a45dd52a3113fff16e6beda69434199))
* add support for passing files as parameters ([5a44115](https://github.com/stainless-api/stainless-api-cli/commit/5a44115f285231ae36c872a9ad510b48b8d4f55c))
* **api:** update support email address ([a9113c7](https://github.com/stainless-api/stainless-api-cli/commit/a9113c7c26bb474e3118a1d4f096763881cafa15))
* **client:** provide file completions when using file embed syntax ([d480ff2](https://github.com/stainless-api/stainless-api-cli/commit/d480ff2c074f611580f0a36cf96f316d5bc516ca))
* **cli:** improve shell completions for namespaced commands and flags ([d11eba5](https://github.com/stainless-api/stainless-api-cli/commit/d11eba51ee4b07e1360cefde3e329c24d47067d0))
* improved support for passing files for `any`-typed arguments ([73f2409](https://github.com/stainless-api/stainless-api-cli/commit/73f24091a362279a62005946079d6c6f301fc4b4))


### Bug Fixes

* cleanup and make passing project flag make more sense ([91a6b99](https://github.com/stainless-api/stainless-api-cli/commit/91a6b9977b85c115ca97eee02e655401809c19e3))
* fix for file uploads to octet stream and form encoding endpoints ([55530da](https://github.com/stainless-api/stainless-api-cli/commit/55530da089f0017c9815c456b5a2f850a241f122))
* fix for nullable arguments ([7bbd057](https://github.com/stainless-api/stainless-api-cli/commit/7bbd05715c474f582ce2d076831d01a6fb32e2e4))
* fix for when terminal width is not available ([78d7919](https://github.com/stainless-api/stainless-api-cli/commit/78d791948a88bb9e509da315f8b51b89d452e945))
* fix mock tests with inner fields that have underscores ([0fb140a](https://github.com/stainless-api/stainless-api-cli/commit/0fb140aa40f3e52e01a9a8b8142876e9fcee8647))
* preserve filename in content-disposition for file uploads ([aae338c](https://github.com/stainless-api/stainless-api-cli/commit/aae338c670ad9938b1b2835ceba3407f25bcb2ba))
* prevent tests from hanging on streaming/paginated endpoints ([257215b](https://github.com/stainless-api/stainless-api-cli/commit/257215b2d77017373910dbc5558e239b7e5a5b7c))
* use RawJSON for iterated values instead of re-marshalling ([fcc55a6](https://github.com/stainless-api/stainless-api-cli/commit/fcc55a618a14bf956aa0af852d66e463e0073fbf))


### Chores

* add build step to ci ([c957b31](https://github.com/stainless-api/stainless-api-cli/commit/c957b31381e935c9f8c80738bf003deb5b05727e))
* update documentation in readme ([dc0fd28](https://github.com/stainless-api/stainless-api-cli/commit/dc0fd28e141c8f50c90a47144d7257d596612ada))


### Refactors

* extract workspaceconfig.go ([23cf851](https://github.com/stainless-api/stainless-api-cli/commit/23cf85125927e0d75b910c7dcb54e300e5b8dc0c))

## 0.1.0-alpha.70 (2026-01-28)

Full Changelog: [v0.1.0-alpha.69...v0.1.0-alpha.70](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.69...v0.1.0-alpha.70)

### Features

* set project to be not required ([83baf59](https://github.com/stainless-api/stainless-api-cli/commit/83baf59ea5f31e616a804e162635074c1c48d8bf))


### Bug Fixes

* **init:** add project flag ([d727066](https://github.com/stainless-api/stainless-api-cli/commit/d727066630bdfd558c435e68d6616e7deb9cbc28))


### Chores

* bump pkg versions ([5a77ff1](https://github.com/stainless-api/stainless-api-cli/commit/5a77ff18cac1ddaedfe1b81166b1db384787ea8e))

## 0.1.0-alpha.69 (2026-01-23)

Full Changelog: [v0.1.0-alpha.68...v0.1.0-alpha.69](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.68...v0.1.0-alpha.69)

### Features

* **api:** ai commit message features are available as public feature ([b8e3cbb](https://github.com/stainless-api/stainless-api-cli/commit/b8e3cbb5b5392955488f3d2821deb937d050e60a))
* **api:** manual updates ([8f23799](https://github.com/stainless-api/stainless-api-cli/commit/8f23799bac4d79fc4221ed698fe3415766a6d9f4))


### Bug Fixes

* don't overwrite flag if revision is explicitly set ([fac264a](https://github.com/stainless-api/stainless-api-cli/commit/fac264a2ab469f61912c8934632e59f7749bebbd))


### Chores

* **internal:** update `actions/checkout` version ([7636b1e](https://github.com/stainless-api/stainless-api-cli/commit/7636b1e3af2dea64143fdb300b90e94d9938ef74))

## 0.1.0-alpha.68 (2026-01-14)

Full Changelog: [v0.1.0-alpha.67...v0.1.0-alpha.68](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.67...v0.1.0-alpha.68)

### Bug Fixes

* avoid consuming request bodies when printing redacted outputs ([861b9b3](https://github.com/stainless-api/stainless-api-cli/commit/861b9b3dadc369ab7ae4c215318ad2102c7be7a4))
* fix terminal height issues causing test failures ([59111de](https://github.com/stainless-api/stainless-api-cli/commit/59111de6cea673462d57d6b82529c1a10b32f449))
* flag defaults ([10e4643](https://github.com/stainless-api/stainless-api-cli/commit/10e46430abd8cb79fb0d5cd8c5f3ec4220c064f6))
* overly broad redaction of Authorization ([6a6dc2e](https://github.com/stainless-api/stainless-api-cli/commit/6a6dc2e53e5957ea017374fc267116310ba8678e))
* prevent flag duplication ([5c9267b](https://github.com/stainless-api/stainless-api-cli/commit/5c9267b768c6453d4cd02e6be23f8e19752b1376))
* stl builds create ([ad0be67](https://github.com/stainless-api/stainless-api-cli/commit/ad0be6704dac6819f50a1d7a3102d618aaf77615))
* stl builds:target-outputs retrieve ([e244860](https://github.com/stainless-api/stainless-api-cli/commit/e244860e9ff0635822642ea84181c6c6c6f831ec))


### Chores

* **deps:** bump golang.org/x/net from 0.33.0 to 0.38.0 ([d38435b](https://github.com/stainless-api/stainless-api-cli/commit/d38435b1c8fcdd52aee5fe606be6c68fa3bcd0eb))
* **internal:** codegen related update ([a31d0d5](https://github.com/stainless-api/stainless-api-cli/commit/a31d0d54ac413e105d3ed869aae8c4b072365d6a))
* update internal comment ([6a6dc2e](https://github.com/stainless-api/stainless-api-cli/commit/6a6dc2e53e5957ea017374fc267116310ba8678e))
* updated README.md with more flag information ([8d8b746](https://github.com/stainless-api/stainless-api-cli/commit/8d8b74622020869858f94fe9eae37b972805bf84))

## 0.1.0-alpha.67 (2026-01-14)

Full Changelog: [v0.1.0-alpha.66...v0.1.0-alpha.67](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.66...v0.1.0-alpha.67)

### Features

* enable suggestion for mistyped commands and flags ([f0bda1c](https://github.com/stainless-api/stainless-api-cli/commit/f0bda1cebc04d7500e29813578ea96043f5a4fb2))


### Bug Fixes

* **client:** do not use pager for short paginated responses ([a67b736](https://github.com/stainless-api/stainless-api-cli/commit/a67b7368fa23656bd7ee2e73db3ee0bda80ba879))
* fix for paginated output not writing to pager correctly ([5664786](https://github.com/stainless-api/stainless-api-cli/commit/566478664ded15ca42585c8306511685047cc2a3))


### Chores

* **internal:** codegen related update ([636fbab](https://github.com/stainless-api/stainless-api-cli/commit/636fbab86afc2cf9a9c7bd3f333356c5bf436dd1))
* **internal:** codegen related update ([67f24aa](https://github.com/stainless-api/stainless-api-cli/commit/67f24aa8e449efcf9092b101e892bc961d8345d2))
* **internal:** codegen related update ([9502a81](https://github.com/stainless-api/stainless-api-cli/commit/9502a81ab539b3eda392cbe3bba2ba53e3462ab1))
* **internal:** codegen related update ([3942483](https://github.com/stainless-api/stainless-api-cli/commit/39424833bdd13c10b1d7c885a97c930ab7daba43))
* **internal:** codegen related update ([a2e2d81](https://github.com/stainless-api/stainless-api-cli/commit/a2e2d811b26af235e6e82e086507d2db590fc8a6))
* **internal:** codegen related update ([bdc6a08](https://github.com/stainless-api/stainless-api-cli/commit/bdc6a0838f0294fc76346ba0ef71416b4cc4d55c))
* **internal:** codegen related update ([600e77b](https://github.com/stainless-api/stainless-api-cli/commit/600e77b2864a93476109d08139fc648740232369))
* update Go SDK version ([8de090d](https://github.com/stainless-api/stainless-api-cli/commit/8de090d46ce5dd634087b7ebc92fcbb9aad7f918))

## 0.1.0-alpha.66 (2026-01-07)

Full Changelog: [v0.1.0-alpha.65...v0.1.0-alpha.66](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.65...v0.1.0-alpha.66)

### Features

* **sql:** initial commit ([c3db520](https://github.com/stainless-api/stainless-api-cli/commit/c3db5205769229bfd17937878a0a63074d73c4cc))

## 0.1.0-alpha.65 (2025-12-22)

Full Changelog: [v0.1.0-alpha.64...v0.1.0-alpha.65](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.64...v0.1.0-alpha.65)

### Features

* added mock server tests ([f62abc2](https://github.com/stainless-api/stainless-api-cli/commit/f62abc2e06e2a0f85e57a22a956e3d9272fd65b2))


### Bug Fixes

* base64 encoding regression ([967d82c](https://github.com/stainless-api/stainless-api-cli/commit/967d82c469ec3f7aca4946c2f3bd4a32f848f7df))
* fix generated flag types and value wrapping ([b8293dc](https://github.com/stainless-api/stainless-api-cli/commit/b8293dcb107347ac57b692d540951b3c09350914))


### Chores

* **cli:** run pre-codegen tests on Windows ([7dd0ffb](https://github.com/stainless-api/stainless-api-cli/commit/7dd0ffba4c671678883eedb41b29dd1816b970da))
* **internal:** codegen related update ([c439ed6](https://github.com/stainless-api/stainless-api-cli/commit/c439ed63774cae542fa6eac8d01095a272061be9))
* **internal:** codegen related update ([f9e9d7d](https://github.com/stainless-api/stainless-api-cli/commit/f9e9d7dbcca54b2df0cde1c84e4bc65f525ef786))

## 0.1.0-alpha.64 (2025-12-17)

Full Changelog: [v0.1.0-alpha.63...v0.1.0-alpha.64](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.63...v0.1.0-alpha.64)

### Bug Fixes

* **cli:** remove `*.exe` files from customer SDK changes ([7809497](https://github.com/stainless-api/stainless-api-cli/commit/780949761c6b49604cb8f18678468f57a38b149c))
* generate flags for parameters specified by client ([51b1031](https://github.com/stainless-api/stainless-api-cli/commit/51b1031279571583416a497fc62f6f3b58e3a0a9))


### Chores

* **cli:** add `*.exe` files back to `.gitignore` ([7882fe1](https://github.com/stainless-api/stainless-api-cli/commit/7882fe18ed7ed8d48d885ce836b437af98863abf))
* **cli:** move `jsonview` subpackage to `internal` ([e3b1c70](https://github.com/stainless-api/stainless-api-cli/commit/e3b1c70a58df206773c0548ed6bc674a835421a0))
* **cli:** temporarily remove `*.exe` from `.gitignore` ([34a0d87](https://github.com/stainless-api/stainless-api-cli/commit/34a0d8706e7498bf08877eab47ccb40d89baf267))
* **internal:** codegen related update ([44b6581](https://github.com/stainless-api/stainless-api-cli/commit/44b6581e5979daf7d1d0a66203134eec1602a9de))
* **internal:** codegen related update ([52c6dc8](https://github.com/stainless-api/stainless-api-cli/commit/52c6dc8d08add4a1b6493790f365f96b86a4ef89))

## 0.1.0-alpha.63 (2025-12-17)

Full Changelog: [v0.1.0-alpha.62...v0.1.0-alpha.63](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.62...v0.1.0-alpha.63)

### Features

* add better suggests when commands don't match ([87295d9](https://github.com/stainless-api/stainless-api-cli/commit/87295d934bb279e38798ad1fb677f88014fd09ff))


### Bug Fixes

* fix merge conflict issues ([5f6cc75](https://github.com/stainless-api/stainless-api-cli/commit/5f6cc757e5ce1227952d2c20a5755d407a75ba35))
* ignore .exe files ([234b806](https://github.com/stainless-api/stainless-api-cli/commit/234b806d58b3223cfd8b25c89ac08529da27a3ee))

## 0.1.0-alpha.62 (2025-12-17)

Full Changelog: [v0.1.0-alpha.61...v0.1.0-alpha.62](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.61...v0.1.0-alpha.62)

### Features

* **ai:** build APIs take SDK level commit messages & new gate for AI commit messages ([76d63c3](https://github.com/stainless-api/stainless-api-cli/commit/76d63c3df26e68aca4cbb670eb8f4b6b409b8111))
* **cli:** automatic streaming for paginated endpoints ([119765f](https://github.com/stainless-api/stainless-api-cli/commit/119765fda59e1af4b94f7604bb34f538c03fb786))
* **cli:** binary request bodies ([b713c9c](https://github.com/stainless-api/stainless-api-cli/commit/b713c9c704550e332fb094fa6c0ad118830a51ef))
* new and simplified CLI flag parsing code and YAML support ([093a401](https://github.com/stainless-api/stainless-api-cli/commit/093a401b9f78950bd73563e9a50fa8fe5ae3bec2))
* new and simplified CLI flag parsing code and YAML support ([bd6e906](https://github.com/stainless-api/stainless-api-cli/commit/bd6e906746a97a8bd0f88e09e22b466f4bdee644))
* redact `Authorization` header when using debug option ([557a499](https://github.com/stainless-api/stainless-api-cli/commit/557a499e428cc15dcb195472fbe3f4c5940deb92))
* redact secrets from other authentication headers when using debug option ([3f38b38](https://github.com/stainless-api/stainless-api-cli/commit/3f38b382557d0dbdcf096f15c7020daf3029a72d))


### Bug Fixes

* **api:** switch 'targets' query param to comma-delimited string in diagnostics endpoint ([8cd21bd](https://github.com/stainless-api/stainless-api-cli/commit/8cd21bd604b9a22265c4ab9853ac9dd39818e486))
* **cli:** fix compilation on Windows ([796e618](https://github.com/stainless-api/stainless-api-cli/commit/796e61843a096fa2b9a893340621919fa789e758))
* fix for empty request bodies ([f417be7](https://github.com/stainless-api/stainless-api-cli/commit/f417be7223cc8df2d93842b4a0f42dc6ae79ad0c))
* fixed manpage generation ([ebf32ec](https://github.com/stainless-api/stainless-api-cli/commit/ebf32ecfc7f90529870e25833fd773e2b310130d))
* **mcp:** correct code tool API endpoint ([227e3c4](https://github.com/stainless-api/stainless-api-cli/commit/227e3c4cd369c6cec342e92cf03c50de0c7e74d5))
* paginated endpoints now behave better with pagers by default ([86d67d2](https://github.com/stainless-api/stainless-api-cli/commit/86d67d2a5a48872b6375a236d50f6a92ed7c6654))


### Chores

* **internal:** codegen related update ([a2daf2f](https://github.com/stainless-api/stainless-api-cli/commit/a2daf2f74f4f07c04eb7be94b12a158727bf265a))
* **internal:** codegen related update ([dd0f6e2](https://github.com/stainless-api/stainless-api-cli/commit/dd0f6e20855e397d99053d5b8755c8589cf1c042))
* **internal:** codegen related update ([983d24c](https://github.com/stainless-api/stainless-api-cli/commit/983d24c0fed111746bbff285bb930bd53751ed29))
* **internal:** codegen related update ([479bceb](https://github.com/stainless-api/stainless-api-cli/commit/479bceb648d93bea9a0132a0c99b84780ef78f8b))
* **internal:** codegen related update ([ff44c29](https://github.com/stainless-api/stainless-api-cli/commit/ff44c298f7e066820360bc2fca539957c4c66c94))
* **internal:** version bump ([040e124](https://github.com/stainless-api/stainless-api-cli/commit/040e1243ea5ba1264ea35cb5039af40bb51c0e4f))
* use `stretchr/testify` assertion helpers in tests ([c399dc9](https://github.com/stainless-api/stainless-api-cli/commit/c399dc9ec460998a23c6701ae4c40924d3fed08e))

## 0.1.0-alpha.61 (2025-12-15)

Full Changelog: [v0.1.0-alpha.60...v0.1.0-alpha.61](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.60...v0.1.0-alpha.61)

## 0.1.0-alpha.60 (2025-12-08)

Full Changelog: [v0.1.0-alpha.59...v0.1.0-alpha.60](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.59...v0.1.0-alpha.60)

### Features

* add spinner for project creation ([361bfec](https://github.com/stainless-api/stainless-api-cli/commit/361bfec53b2cdc963e7599626ffdae65f505f8fb))
* fix edge cases for sending request data and add YAML support ([99caa6a](https://github.com/stainless-api/stainless-api-cli/commit/99caa6af3ae82d5622de085066c370cee3e73c71))


### Bug Fixes

* fix for default flag values ([114f9c1](https://github.com/stainless-api/stainless-api-cli/commit/114f9c120eac20cc5a24e2bc48a0049c268d20a9))


### Chores

* **internal:** codegen related update ([858b183](https://github.com/stainless-api/stainless-api-cli/commit/858b183ae8b9cd215876db334b4a6992d7403d47))

## 0.1.0-alpha.59 (2025-12-04)

Full Changelog: [v0.1.0-alpha.58...v0.1.0-alpha.59](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.58...v0.1.0-alpha.59)

### Features

* checkout branch if we can ([a46bd17](https://github.com/stainless-api/stainless-api-cli/commit/a46bd17a96ff604e5211f416d4dad26912b065ad))
* fix builds:target-outputs --pull ([0b4fe3f](https://github.com/stainless-api/stainless-api-cli/commit/0b4fe3f27768af5b2ce3bdcea97e971976230e24))
* fix builds:target-outputs --pull ([3768058](https://github.com/stainless-api/stainless-api-cli/commit/3768058960e3e097ca4db290dc25663e28e4b895))
* handle merge conflicts better ([80c423f](https://github.com/stainless-api/stainless-api-cli/commit/80c423f9b7ffd69c877a2e21f61d5928294f693a))
* show error if download fails ([53df302](https://github.com/stainless-api/stainless-api-cli/commit/53df302d843b5ff34e39a6702ca83faf582e18e8))


### Bug Fixes

* codegen merge conflict issues ([3c557b6](https://github.com/stainless-api/stainless-api-cli/commit/3c557b6f210a4d4474545a88ce1f90f8d8aa8307))
* handle edge cases in build view component better ([a3bdf09](https://github.com/stainless-api/stainless-api-cli/commit/a3bdf0979aff0676616e4ec5dce8e833ad4d72d7))
* improve file paths rendering ([8470724](https://github.com/stainless-api/stainless-api-cli/commit/847072420f6d497d8e9f25396d4bd3106a463033))


### Chores

* attempt to fix merge conflicts ([3aaa391](https://github.com/stainless-api/stainless-api-cli/commit/3aaa3910313a0c2cdb2c3cb9578cb87962e68df3))
* format ([809894d](https://github.com/stainless-api/stainless-api-cli/commit/809894d90ecdd500f2182fd0231cabf9619f820d))
* **internal:** codegen related update ([e61afe2](https://github.com/stainless-api/stainless-api-cli/commit/e61afe2c51d989df5c69296045a0fa883c8e29ee))
* **internal:** codegen related update ([cf7e344](https://github.com/stainless-api/stainless-api-cli/commit/cf7e34445270b8da6608559d24f5e2431e7752a0))
* skip spec resource ([c677113](https://github.com/stainless-api/stainless-api-cli/commit/c6771132483d0f06ea7f8f887b7c9825b9fb61af))
* update dependencies ([e4839cf](https://github.com/stainless-api/stainless-api-cli/commit/e4839cfc6330f355906652a35f3c4047e9d94093))

## 0.1.0-alpha.58 (2025-11-26)

Full Changelog: [v0.1.0-alpha.57...v0.1.0-alpha.58](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.57...v0.1.0-alpha.58)

### Bug Fixes

* consider merge conflict done ([f9ad68b](https://github.com/stainless-api/stainless-api-cli/commit/f9ad68b2a21c9f73f03aa6260f929288e94e680c))
* fix auth issue ([11d5375](https://github.com/stainless-api/stainless-api-cli/commit/11d5375223687be5951f2f1eb3d77035e4ca00ba))


### Refactors

* rename commit step to codegen ([58e752a](https://github.com/stainless-api/stainless-api-cli/commit/58e752a40328b2e8c1d60ee3576daa92b59cc6f2))

## 0.1.0-alpha.57 (2025-11-25)

Full Changelog: [v0.1.0-alpha.56...v0.1.0-alpha.57](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.56...v0.1.0-alpha.57)

### Bug Fixes

* fix init auth ([c840610](https://github.com/stainless-api/stainless-api-cli/commit/c84061006365d4236bab054b335ad596cf18d873))

## 0.1.0-alpha.56 (2025-11-25)

Full Changelog: [v0.1.0-alpha.55...v0.1.0-alpha.56](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.55...v0.1.0-alpha.56)

### Features

* convert paths to absolute path in the .stainless/workspace.json ([14aeb23](https://github.com/stainless-api/stainless-api-cli/commit/14aeb23dbac9e1a6a3f24d88a50e349eef281502))


### Refactors

* significantly refactor builds viewer ([0677658](https://github.com/stainless-api/stainless-api-cli/commit/067765817d047937000f634731bb8aeea6e60346))

## 0.1.0-alpha.55 (2025-11-17)

Full Changelog: [v0.1.0-alpha.54...v0.1.0-alpha.55](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.54...v0.1.0-alpha.55)

### Features

* clean up messages / formatting ([eb3b9b1](https://github.com/stainless-api/stainless-api-cli/commit/eb3b9b13d804ce10327bcd35ce03704f36d77c68))


### Chores

* move dev_view ([0512c7f](https://github.com/stainless-api/stainless-api-cli/commit/0512c7f57dc27d0461bcdce7da9bee346bf6e052))


### Refactors

* extract console ([83f3da7](https://github.com/stainless-api/stainless-api-cli/commit/83f3da7ffeed20d1b9dbb4c9bb23b771ac7f9234))
* extract stainlessutils ([9ee2138](https://github.com/stainless-api/stainless-api-cli/commit/9ee21382a58eaa4ebee33192524c909bf341e237))
* extract stainlessviews ([d075dea](https://github.com/stainless-api/stainless-api-cli/commit/d075dea1085b36d5761eea19253cb253c5273eca))

## 0.1.0-alpha.54 (2025-11-14)

Full Changelog: [v0.1.0-alpha.53...v0.1.0-alpha.54](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.53...v0.1.0-alpha.54)

## 0.1.0-alpha.53 (2025-11-13)

Full Changelog: [v0.1.0-alpha.51...v0.1.0-alpha.53](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.51...v0.1.0-alpha.53)

### Features

* **api:** add branch reset functionality ([b0df2a1](https://github.com/stainless-api/stainless-api-cli/commit/b0df2a1520a8c57fb854e8dc3b133ce195c122d5))


### Chores

* bump go sdk version ([d34a873](https://github.com/stainless-api/stainless-api-cli/commit/d34a87386da62d51ef9fd33fd5637de54d11ee67))
* **internal:** codegen related update ([4e786ec](https://github.com/stainless-api/stainless-api-cli/commit/4e786ece778e52ced853f83374e2fefd3c8460d1))

## 0.1.0-alpha.51 (2025-10-25)

Full Changelog: [v0.1.0-alpha.50...v0.1.0-alpha.51](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.50...v0.1.0-alpha.51)

### Features

* use scope:"*" ([6b45c41](https://github.com/stainless-api/stainless-api-cli/commit/6b45c41e78d6992355533234986ee07f212d92c0))


### Bug Fixes

* fix builds for non-public Go repos ([725ff05](https://github.com/stainless-api/stainless-api-cli/commit/725ff05a2078bef1e4081e8a19bc3be12ce20a0f))
* remove some bootstrapping logic ([a3f713d](https://github.com/stainless-api/stainless-api-cli/commit/a3f713d5551ca6c0b6d1bdf4d83999652c15d9cf))

## 0.1.0-alpha.50 (2025-10-23)

Full Changelog: [v0.1.0-alpha.49...v0.1.0-alpha.50](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.49...v0.1.0-alpha.50)

### Features

* arguments now have defaults and descriptions ([9b8d534](https://github.com/stainless-api/stainless-api-cli/commit/9b8d534c8b92ae4f03dda9ac12e80c0574a0642d))
* Expose connection-specific decorated OAS ([abf1d87](https://github.com/stainless-api/stainless-api-cli/commit/abf1d87e0af1e8f77405873c1ddf36f5070f897c))


### Bug Fixes

* pass through context parameter correctly ([29182fd](https://github.com/stainless-api/stainless-api-cli/commit/29182fdae10f5f91e3b8d0f90234f4e330572c9d))


### Chores

* bump go sdk version ([358eed4](https://github.com/stainless-api/stainless-api-cli/commit/358eed41b58f7bacf6d5ce340795d251fa1803c3))
* bump Go version ([84c6127](https://github.com/stainless-api/stainless-api-cli/commit/84c612747d1a22ee69d3ffca48f7c538b9b44896))
* **internal:** codegen related update ([375ca6b](https://github.com/stainless-api/stainless-api-cli/commit/375ca6ba3c64a74577430abec867b1b9b7c2fe43))

## 0.1.0-alpha.49 (2025-10-06)

Full Changelog: [v0.1.0-alpha.48...v0.1.0-alpha.49](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.48...v0.1.0-alpha.49)

### Bug Fixes

* downgrade urfave/cli-docs dependency ([6b4c1ac](https://github.com/stainless-api/stainless-api-cli/commit/6b4c1accf900e7b7e6266873fa46e7153b15e7ec))

## 0.1.0-alpha.48 (2025-10-03)

Full Changelog: [v0.1.0-alpha.47...v0.1.0-alpha.48](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.47...v0.1.0-alpha.48)

## 0.1.0-alpha.47 (2025-10-03)

Full Changelog: [v0.1.0-alpha.46...v0.1.0-alpha.47](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.46...v0.1.0-alpha.47)

## 0.1.0-alpha.46 (2025-10-03)

Full Changelog: [v0.1.0-alpha.45...v0.1.0-alpha.46](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.45...v0.1.0-alpha.46)

## 0.1.0-alpha.45 (2025-10-02)

Full Changelog: [v0.1.0-alpha.44...v0.1.0-alpha.45](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.44...v0.1.0-alpha.45)

### Features

* added `--output-filter` flag and `--error-format` flag to support better visualization options ([2768e5b](https://github.com/stainless-api/stainless-api-cli/commit/2768e5b0255ae79e44ed54f583c4809f19801afe))
* better support for positional arguments ([799d88f](https://github.com/stainless-api/stainless-api-cli/commit/799d88f441c5d1d5dca82a5590fceefdc1823f42))


### Chores

* **internal:** codegen related update ([2f9d764](https://github.com/stainless-api/stainless-api-cli/commit/2f9d764f660458d0e84e11adea9f1d418e8f1013))
* **internal:** codegen related update ([c7df1a7](https://github.com/stainless-api/stainless-api-cli/commit/c7df1a7401a51fe9657ca9fa7f1aff2d37f60851))

## 0.1.0-alpha.44 (2025-09-22)

Full Changelog: [v0.1.0-alpha.43...v0.1.0-alpha.44](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.43...v0.1.0-alpha.44)

### Features

* **api:** docs ([f0762b0](https://github.com/stainless-api/stainless-api-cli/commit/f0762b08842d63a55f28afc18a77907c24f2513b))
* improved formatting options for command outputs ([f5e6481](https://github.com/stainless-api/stainless-api-cli/commit/f5e6481706b882eaa0c84a3bf8aea2a088087402))
* show full error message on fatal error ([00621a1](https://github.com/stainless-api/stainless-api-cli/commit/00621a1ba20a392670dc8a7bf991c1d0507efc12))


### Bug Fixes

* fix for issue with nil responses ([14f22bf](https://github.com/stainless-api/stainless-api-cli/commit/14f22bfbe6318e849a57b2243c958e3af869b18c))
* fix go client version bump issues ([d59e6ea](https://github.com/stainless-api/stainless-api-cli/commit/d59e6ea0f8363c1ad7ace9c97e243ec972f75c24))


### Chores

* code cleanup for `interface{}` ([30cd6f1](https://github.com/stainless-api/stainless-api-cli/commit/30cd6f10ce33163f4f6b7f13085a259cf738ca49))
* do not install brew dependencies in ./scripts/bootstrap by default ([ff84f98](https://github.com/stainless-api/stainless-api-cli/commit/ff84f986702f63a60e1f5343e180c305877f6171))
* update go dependency ([22a15b8](https://github.com/stainless-api/stainless-api-cli/commit/22a15b85720596afe3dd95d3be4f2ce4ddf1341a))

## 0.1.0-alpha.43 (2025-09-16)

Full Changelog: [v0.1.0-alpha.42...v0.1.0-alpha.43](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.42...v0.1.0-alpha.43)

### Bug Fixes

* fix diagnostic iteration ([bd8ce45](https://github.com/stainless-api/stainless-api-cli/commit/bd8ce45278e05dcc90ee88bdb9dfb9225345c0f9))

## 0.1.0-alpha.42 (2025-09-15)

Full Changelog: [v0.1.0-alpha.41...v0.1.0-alpha.42](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.41...v0.1.0-alpha.42)

### Features

* improved formatting options for command outputs ([0fef47d](https://github.com/stainless-api/stainless-api-cli/commit/0fef47d958d79dea5d71674bf6eee5109a14743d))
* skip ignored diagnostics ([1edf832](https://github.com/stainless-api/stainless-api-cli/commit/1edf8322fe459c6c2fdf524e3501e5efd03bc97e))


### Bug Fixes

* dont crash when git user is not available ([7012171](https://github.com/stainless-api/stainless-api-cli/commit/701217124580a992cc47fdd8c4762f70c1130105))
* lint error ([7eda049](https://github.com/stainless-api/stainless-api-cli/commit/7eda0495e88be8348edf9c789ce60f0cac6eab3f))
* more merge conflict issues ([2416212](https://github.com/stainless-api/stainless-api-cli/commit/2416212db96b9f3fb15edc5c08ff7a9af8bb1d0f))

## 0.1.0-alpha.41 (2025-09-09)

Full Changelog: [v0.1.0-alpha.40...v0.1.0-alpha.41](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.40...v0.1.0-alpha.41)

### Bug Fixes

* **homebrew:** homebrew distribution should work now ([db54fbd](https://github.com/stainless-api/stainless-api-cli/commit/db54fbd09b4b259e054300a5d711f5798e425cde))

## 0.1.0-alpha.40 (2025-09-09)

Full Changelog: [v0.1.0-alpha.39...v0.1.0-alpha.40](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.39...v0.1.0-alpha.40)

### Features

* **api:** resources and methods now use kebab-case instead of snake_case ([30eb64b](https://github.com/stainless-api/stainless-api-cli/commit/30eb64bd5450582b8a15d12720550574f9f6567e))
* **auth:** add descriptive device names with hostname and OS info ([b36e2d0](https://github.com/stainless-api/stainless-api-cli/commit/b36e2d02d7f34b97aa44cea21993e74791d4d43e))
* now ships with manpages ([1f312b4](https://github.com/stainless-api/stainless-api-cli/commit/1f312b4930b5245866fd8fce1f5266962f7ee65d))


### Bug Fixes

* some methods no longer require a prefix ([4251212](https://github.com/stainless-api/stainless-api-cli/commit/425121255ec5f7abcb909a919fa0b20bc1ab18f9))


### Chores

* bump go sdk version ([9669c15](https://github.com/stainless-api/stainless-api-cli/commit/9669c158faf4f88861b2c4016c54d453d42c698f))
* **internal:** codegen related update ([291d111](https://github.com/stainless-api/stainless-api-cli/commit/291d111fa1512ddb47cab59d94a0524b440f4ed2))

## 0.1.0-alpha.39 (2025-08-25)

Full Changelog: [v0.1.0-alpha.38...v0.1.0-alpha.39](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.38...v0.1.0-alpha.39)

### Features

* **api:** manual updates ([351b067](https://github.com/stainless-api/stainless-api-cli/commit/351b067c40ae8dde9369a6557b1a98b42dde835e))


### Documentation

* add contact email and link to docs ([f382570](https://github.com/stainless-api/stainless-api-cli/commit/f382570eb1016f346ada5fa64a0c668a2b01b52d))

## 0.1.0-alpha.38 (2025-08-20)

Full Changelog: [v0.1.0-alpha.37...v0.1.0-alpha.38](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.37...v0.1.0-alpha.38)

### Features

* [WIP] add branch rebasing to the API ([0ffdf30](https://github.com/stainless-api/stainless-api-cli/commit/0ffdf303808ea76b14c80280ab36d94a9f6d739d))
* **api:** update go_sdk_version ([dcb1497](https://github.com/stainless-api/stainless-api-cli/commit/dcb14971ba57548535646186d3ccd1a0ec8830d1))
* make --branch no longer required ([ecea04a](https://github.com/stainless-api/stainless-api-cli/commit/ecea04aae24b722fd046fa80a39f380157e7e3e0))

## 0.1.0-alpha.37 (2025-08-13)

Full Changelog: [v0.1.0-alpha.36...v0.1.0-alpha.37](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.36...v0.1.0-alpha.37)

### Features

* add User-Agent ([5e14044](https://github.com/stainless-api/stainless-api-cli/commit/5e140445516ab732da0290c46f6d185df933d258))
* infer branch better ([06eb8ab](https://github.com/stainless-api/stainless-api-cli/commit/06eb8ab1098dc2fc6d66bf44021a5e5ed4f5217f))


### Chores

* internal change ([7acb27b](https://github.com/stainless-api/stainless-api-cli/commit/7acb27b713b5530d017ef8fa1f04c1edaf8bb636))
* **internal:** codegen related update ([31f37f8](https://github.com/stainless-api/stainless-api-cli/commit/31f37f88464ef707f93d32b1dbfa10584b4c12b9))
* **internal:** codegen related update ([9f7e69f](https://github.com/stainless-api/stainless-api-cli/commit/9f7e69fcd733ccfa9e37bd09b599edaf4573f7d9))
* update @stainless-api/prism-cli to v5.15.0 ([555a3e3](https://github.com/stainless-api/stainless-api-cli/commit/555a3e346fb4d5eee9d8b16d6769bdbc49b0be9a))


### Refactors

* extract logic out to stainlessutils.go ([a23a33a](https://github.com/stainless-api/stainless-api-cli/commit/a23a33aaa08129f6a00d0423d2660d91aaaf6d3a))
* extract logic to authconfig.go ([c140592](https://github.com/stainless-api/stainless-api-cli/commit/c14059238b6e9c43fbc55f98271d7ae608a8e6d0))
* extract logic to workspaceconfig.go ([8735957](https://github.com/stainless-api/stainless-api-cli/commit/87359570d62e6567d2d9c9288e82af22fecd6142))
* simplify AuthConfig interface ([27572ad](https://github.com/stainless-api/stainless-api-cli/commit/27572ad0621b45d610978083880f2a46a9388378))

## 0.1.0-alpha.36 (2025-08-06)

Full Changelog: [v0.1.0-alpha.35...v0.1.0-alpha.36](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.35...v0.1.0-alpha.36)

### Features

* make build resolve much faster ([922d86c](https://github.com/stainless-api/stainless-api-cli/commit/922d86c67c6c28a98859f5cf4b70eca9eddf9f39))
* skip downloads for builds with fatal commit conclusions ([6cb7a29](https://github.com/stainless-api/stainless-api-cli/commit/6cb7a299fc24ca809b540489bda0d9d2e0ceb3ac))
* skip downloads in dev mode for fatal commit conclusions ([228bd08](https://github.com/stainless-api/stainless-api-cli/commit/228bd082f08848fa40e547677f31d70ba469e812))


### Bug Fixes

* don't set openapi-spec if revision is set ([ddc3715](https://github.com/stainless-api/stainless-api-cli/commit/ddc37155f67018d8a3839acfe08aa4772b3040a4))


### Chores

* improve messaging when stainless config isn't configured ([50b5348](https://github.com/stainless-api/stainless-api-cli/commit/50b5348ccc13415e07293e07ab5c73168a2c03b3))


### Refactors

* remove retrying from config download ([af51b3c](https://github.com/stainless-api/stainless-api-cli/commit/af51b3c21e07419b3d3044a3a7d590a31c131815))

## 0.1.0-alpha.35 (2025-08-05)

Full Changelog: [v0.1.0-alpha.34...v0.1.0-alpha.35](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.34...v0.1.0-alpha.35)

### Features

* add path to config output ([7a6fd0b](https://github.com/stainless-api/stainless-api-cli/commit/7a6fd0bf6f227b4c9c79b92294882130d9539e03))
* add remote if it doesn't exist ([8455a8d](https://github.com/stainless-api/stainless-api-cli/commit/8455a8d2c715da64c2041c37854200fe202ac075))
* change next steps padding ([e2e3be2](https://github.com/stainless-api/stainless-api-cli/commit/e2e3be2e211fc0cf7503529041f2220566ea0de0))
* fix workspace init command ([f7873f3](https://github.com/stainless-api/stainless-api-cli/commit/f7873f3c3c43508558a58b8d4e987aaa634316d4))
* improve form styling ([290d3a6](https://github.com/stainless-api/stainless-api-cli/commit/290d3a62b810be711b28f8b43e94a6dd52cd4191))
* remove too much choice from init flow ([056aeec](https://github.com/stainless-api/stainless-api-cli/commit/056aeec148d3ba28748682b86d47669967a4ba7f))
* strip http login url ([6576c07](https://github.com/stainless-api/stainless-api-cli/commit/6576c07f9b85dc2b1c35eca392ac09540a97b2b9))
* use checkmark for multiselect ([3edf4be](https://github.com/stainless-api/stainless-api-cli/commit/3edf4bea0a14607b22768f821aa558941ec47353))


### Chores

* add usage examples ([2f51c69](https://github.com/stainless-api/stainless-api-cli/commit/2f51c699a2ed4ad8ebff9137e8687f6f19648d2e))
* consistently capitalize Stainless ([9f3d5e8](https://github.com/stainless-api/stainless-api-cli/commit/9f3d5e829ef18b9977dfd8bbf88ee710ec19ec82))

## 0.1.0-alpha.34 (2025-07-31)

Full Changelog: [v0.1.0-alpha.33...v0.1.0-alpha.34](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.33...v0.1.0-alpha.34)

### Features

* **api:** manual updates ([ca8a538](https://github.com/stainless-api/stainless-api-cli/commit/ca8a538790ae9599f9eb0ec07b8c04c377405751))
* simplify getAPICommandContext ([23ccc6d](https://github.com/stainless-api/stainless-api-cli/commit/23ccc6db032a8b7bba5b86bf1ecb369c15fba7f7))
* wait for organization creation in CLI ([2295a90](https://github.com/stainless-api/stainless-api-cli/commit/2295a90f35b25e7e30dccb457024ba5fdfa0f271))


### Chores

* update tap repository ([7a8cb8d](https://github.com/stainless-api/stainless-api-cli/commit/7a8cb8d175b94c7fd8e017665769a359d3f52f8f))

## 0.1.0-alpha.33 (2025-07-30)

Full Changelog: [v0.1.0-alpha.32...v0.1.0-alpha.33](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.32...v0.1.0-alpha.33)

### Bug Fixes

* fix git fetch when building branch ([7cb800e](https://github.com/stainless-api/stainless-api-cli/commit/7cb800ef4717ddbf39296f239a31678ab3251540))

## 0.1.0-alpha.32 (2025-07-30)

Full Changelog: [v0.1.0-alpha.31...v0.1.0-alpha.32](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.31...v0.1.0-alpha.32)

### Features

* add download to dev mode ([6532712](https://github.com/stainless-api/stainless-api-cli/commit/65327127654bf67a804aeb141cfec4a4017c183c))


### Bug Fixes

* remove debug logging ([87d0f51](https://github.com/stainless-api/stainless-api-cli/commit/87d0f51c63f89ca42dcd4731ca872faeecf30e01))

## 0.1.0-alpha.31 (2025-07-30)

Full Changelog: [v0.1.0-alpha.30...v0.1.0-alpha.31](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.30...v0.1.0-alpha.31)

### Features

* improve autofill project name logic ([6adcbcd](https://github.com/stainless-api/stainless-api-cli/commit/6adcbcddc9679f89c3da93d23e8315a8fb6437be))


### Bug Fixes

* remove unused os import ([a8d8227](https://github.com/stainless-api/stainless-api-cli/commit/a8d8227e64f46c5237da2a4dedcebab389d8ad02))


### Refactors

* clean up file flags ([7248111](https://github.com/stainless-api/stainless-api-cli/commit/72481118d3fc5620bd91dfea8e57849c97caffde))

## 0.1.0-alpha.30 (2025-07-29)

Full Changelog: [v0.1.0-alpha.29...v0.1.0-alpha.30](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.29...v0.1.0-alpha.30)

### Features

* add links to dev mode ([579819a](https://github.com/stainless-api/stainless-api-cli/commit/579819a364d8527000472c47923f64efc0dff8d9))
* don't print project json on init ([f5479c4](https://github.com/stainless-api/stainless-api-cli/commit/f5479c4949ce837ac1cc1b92d4089a5bb15af363))
* fix spacing between waiting for build and pulling outputs ([472cb40](https://github.com/stainless-api/stainless-api-cli/commit/472cb40d9ffe818d964e3b82ec44e8ec56c7379a))
* support reading from config.Targets in builds create ([a451f8f](https://github.com/stainless-api/stainless-api-cli/commit/a451f8f1f59588bca5c0b00b4081976b2aa65a59))
* update target download logic ([3d4a32e](https://github.com/stainless-api/stainless-api-cli/commit/3d4a32e17fee21a647845db2edf52892ecf80d26))


### Refactors

* change semantics of downloadStainlessConfig ([2748f4b](https://github.com/stainless-api/stainless-api-cli/commit/2748f4bfd4e7c09075046abb11e7c8ec4021622b))
* simplify config threading ([3e3874a](https://github.com/stainless-api/stainless-api-cli/commit/3e3874aa6536e5243dc6409b27ac487dd38a4957))

## 0.1.0-alpha.29 (2025-07-29)

Full Changelog: [v0.1.0-alpha.28...v0.1.0-alpha.29](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.28...v0.1.0-alpha.29)

### Features

* add diagnositcs count ([fa36fb2](https://github.com/stainless-api/stainless-api-cli/commit/fa36fb23190f91a79c2547bcad1def4f66497c69))
* add init command ([af9d04b](https://github.com/stainless-api/stainless-api-cli/commit/af9d04b6aae5d4b97df29c8f48257fd26377d82e))
* add targets to workspace config ([49e41a2](https://github.com/stainless-api/stainless-api-cli/commit/49e41a24b8f574df8f1e7641329192e1af43a20d))
* port some init logic to workspace init ([c509357](https://github.com/stainless-api/stainless-api-cli/commit/c509357a229e08a2d470db7e664cfa1248903191))
* pull targets from init ([98d427b](https://github.com/stainless-api/stainless-api-cli/commit/98d427b06f1b0509e2ea5745184d27863119c5a2))
* use getAvailableTargetInfo for preselected builds ([63b7b1d](https://github.com/stainless-api/stainless-api-cli/commit/63b7b1d7d5090deeba60efbef461fb6f87dd2416))


### Chores

* run format ([f2703f4](https://github.com/stainless-api/stainless-api-cli/commit/f2703f41c1a8a662fbfa615ee582e9040c0b4e8c))


### Refactors

* rename getCompletedTargets to getTargetInfo ([a44754a](https://github.com/stainless-api/stainless-api-cli/commit/a44754a8b14e03b51716ec0e9d8f6baa948da43f))

## 0.1.0-alpha.28 (2025-07-28)

Full Changelog: [v0.1.0-alpha.27...v0.1.0-alpha.28](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.27...v0.1.0-alpha.28)

### Features

* add callout at project creation ([003796e](https://github.com/stainless-api/stainless-api-cli/commit/003796e6004bc9c9986861a72a4288ab393aa953))
* add some polish around confirm dialogs ([40d3419](https://github.com/stainless-api/stainless-api-cli/commit/40d34192c6cfc937d9da9ddb4f9b393bf808417c))
* **api:** manual updates ([27041d4](https://github.com/stainless-api/stainless-api-cli/commit/27041d4a81fe0cee98907a872bd73b0c60f92b31))
* **api:** manual updates ([b97d699](https://github.com/stainless-api/stainless-api-cli/commit/b97d699c3e1d17dc7199ffdf8768c4942d875936))
* flesh out project branches endpoints ([897282f](https://github.com/stainless-api/stainless-api-cli/commit/897282f565bc519c4256aef101aadb36c16e6c96))
* improve output grouping ([012926f](https://github.com/stainless-api/stainless-api-cli/commit/012926fc3c5b0e44243a689e4fd911baebe46395))
* improve project create form ([e3ec063](https://github.com/stainless-api/stainless-api-cli/commit/e3ec0632a86ae39d58f056a187636d3fc06c6c24))
* remove stainless config uploading to the project create endpoint ([3bba950](https://github.com/stainless-api/stainless-api-cli/commit/3bba950d20c4117103f0ff1cb936f80543f80c0e))
* rename init-workspace flag to workspace-init ([c3c7887](https://github.com/stainless-api/stainless-api-cli/commit/c3c7887e2bc39a3c7b20f58989f5cb7eab4ebcc4))


### Bug Fixes

* add retries to project config retrieve ([630b8e8](https://github.com/stainless-api/stainless-api-cli/commit/630b8e82c970478cc481093c9de774f08d55ff9b))
* don't unconditionally indent forms ([249b151](https://github.com/stainless-api/stainless-api-cli/commit/249b15100bb714886efcfa8093b2108cda3fa5e4))
* fix middleware with empty body ([b734f2e](https://github.com/stainless-api/stainless-api-cli/commit/b734f2ee032645f064f4c033131c164c75be6e97))
* improve diagnostics printing ([9d42891](https://github.com/stainless-api/stainless-api-cli/commit/9d42891284b59e3c976f9eae879e114539c2ff9e))
* improve printing of diagnostics when no diagnostics are there ([65e18d5](https://github.com/stainless-api/stainless-api-cli/commit/65e18d5c1d4167a4ae5911df3ed980f064c5c57e))
* unwrap content for stainless config ([919d1b7](https://github.com/stainless-api/stainless-api-cli/commit/919d1b72fd7b8c9482039466ae56ef66efbf61f1))

## 0.1.0-alpha.27 (2025-07-28)

Full Changelog: [v0.1.0-alpha.26...v0.1.0-alpha.27](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.26...v0.1.0-alpha.27)

### Features

* flesh out project branches endpoints ([3f50b90](https://github.com/stainless-api/stainless-api-cli/commit/3f50b90035c92641cebc44cf439fe70695c61098))


### Chores

* bump go version ([6179d9d](https://github.com/stainless-api/stainless-api-cli/commit/6179d9d015abe133f094fb6deca33c70f537f2dd))
* **internal:** codegen related update ([71190a6](https://github.com/stainless-api/stainless-api-cli/commit/71190a66478187d4fdcb129cae63c408f1776631))
* **internal:** codegen related update ([eacc95e](https://github.com/stainless-api/stainless-api-cli/commit/eacc95e59b1cd5cee96ca0567a3e0bb10e1d16cc))
* **internal:** codegen related update ([66c2750](https://github.com/stainless-api/stainless-api-cli/commit/66c275001290e765feea3a636625486f67705380))

## 0.1.0-alpha.26 (2025-07-22)

Full Changelog: [v0.1.0-alpha.25...v0.1.0-alpha.26](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.25...v0.1.0-alpha.26)

### Features

* add dev mode ([5c7d002](https://github.com/stainless-api/stainless-api-cli/commit/5c7d002d71cc1f1955bbd23302e7179e7bfd8658))
* **api:** manual updates ([851173a](https://github.com/stainless-api/stainless-api-cli/commit/851173a3796f1c925b0fcc02e8973b50cacd7dba))
* **api:** manual updates ([f5f13f9](https://github.com/stainless-api/stainless-api-cli/commit/f5f13f97699ec60e69dc8e4f8104f9a154661a0b))
* Codegen(go): encode path params ([3ff985f](https://github.com/stainless-api/stainless-api-cli/commit/3ff985f2f214b5c2a34482dcacbaf796a768ca77))
* Codegen(php): unset optional parameters in constructor ([dcd415a](https://github.com/stainless-api/stainless-api-cli/commit/dcd415a35e1223572e32aae9af7de1b6564046bd))
* php: generate stub union classes with discrimminator info ([dbc8aa0](https://github.com/stainless-api/stainless-api-cli/commit/dbc8aa04f616b9135ac720bc22983271f5db3d2e))


### Chores

* add unstable call out ([a0157be](https://github.com/stainless-api/stainless-api-cli/commit/a0157bef8ca29fc418bc43dd38c453d7ad21bdb8))
* bump go sdk ([c671b1d](https://github.com/stainless-api/stainless-api-cli/commit/c671b1d3ebd95127cccdf97f6b6d5de570c54849))
* bump go sdk version ([6b1aa80](https://github.com/stainless-api/stainless-api-cli/commit/6b1aa80ce7f247841a924adef77207514ebaf18a))
* bump go version ([ad830cf](https://github.com/stainless-api/stainless-api-cli/commit/ad830cf2b99969d0638d1d319f8bac1586dfa8dd))
* update homebrew repository ([e49e98c](https://github.com/stainless-api/stainless-api-cli/commit/e49e98c10098c08f46537cf3c2650d33580d6844))


### Refactors

* improve scrollback content viewer ([140f41b](https://github.com/stainless-api/stainless-api-cli/commit/140f41bd0ed944fe25759ad7bc39643ebfc0e245))
* optimize pipeline rendering ([2ddf648](https://github.com/stainless-api/stainless-api-cli/commit/2ddf64830a410a66b686a082330288fc3683f51e))
* simplify getStepSymbol ([a3f7d60](https://github.com/stainless-api/stainless-api-cli/commit/a3f7d60ea9d32bb73032dcd67673cd00c0d36c8f))

## 0.1.0-alpha.25 (2025-07-15)

Full Changelog: [v0.1.0-alpha.24...v0.1.0-alpha.25](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.24...v0.1.0-alpha.25)

### Features

* add multipart form support to HAR snippets ([24979f2](https://github.com/stainless-api/stainless-api-cli/commit/24979f22f37f1f4e2f6cf473746ce523914f704a))
* add timestamp to builds api ([1144d20](https://github.com/stainless-api/stainless-api-cli/commit/1144d20daa405a332c695fa08087fb43c37c8067))
* **api:** add staging environment ([be9cb34](https://github.com/stainless-api/stainless-api-cli/commit/be9cb344a2ac600947dc4c523a51b5655203cf57))
* **api:** manual updates ([13e06db](https://github.com/stainless-api/stainless-api-cli/commit/13e06db339629317a88a9ae994f0adbfbad49d05))
* **api:** manual updates ([57ee9a1](https://github.com/stainless-api/stainless-api-cli/commit/57ee9a1481625046ca7a9d37bd3ba7f553ac3ca8))
* **api:** manual updates ([5753cc0](https://github.com/stainless-api/stainless-api-cli/commit/5753cc08e608479ba81fbecee3ec07eba252274a))
* improve output coloring ([b653f1c](https://github.com/stainless-api/stainless-api-cli/commit/b653f1c2609edd3f2de16988b995a434ce3fd89b))
* improve output coloring ([617c876](https://github.com/stainless-api/stainless-api-cli/commit/617c876068015ad54412f19e6dd593d126ce5283))
* make build api return documented specs as urls ([1127d78](https://github.com/stainless-api/stainless-api-cli/commit/1127d78574ed81939c6d003372715a5f9a06da47))
* support target:path syntax for custom output directories ([d8c7e4c](https://github.com/stainless-api/stainless-api-cli/commit/d8c7e4ce9e3d7817d8bda3fab1823835122b7f68))
* support target:path syntax for custom output directories ([97e7fae](https://github.com/stainless-api/stainless-api-cli/commit/97e7fae3a5c363c54d5add25c6e72d93a867e2af))
* support target:path syntax for custom output directories ([d5447aa](https://github.com/stainless-api/stainless-api-cli/commit/d5447aaf314ce7b3bb009fe2551d74cc251c3b21))


### Chores

* bump ([e6ba380](https://github.com/stainless-api/stainless-api-cli/commit/e6ba38062edf853d06f9ae47842ebb0137b83240))
* bump go sdk version ([bb75e98](https://github.com/stainless-api/stainless-api-cli/commit/bb75e986444393f502c552b794eba65a295c6dc7))
* **internal:** codegen related update ([e214b20](https://github.com/stainless-api/stainless-api-cli/commit/e214b2087a80dea5a66872a1272713eacd44a2ba))
* **internal:** codegen related update ([8587cc5](https://github.com/stainless-api/stainless-api-cli/commit/8587cc5e08b5887eba04bdd7260de958d7c24943))
* move sdkjson generation api out of v0 scope ([0c7361f](https://github.com/stainless-api/stainless-api-cli/commit/0c7361f9a4e00af8ae6d65a90c1a34cffd586a38))

## 0.1.0-alpha.24 (2025-07-02)

Full Changelog: [v0.1.0-alpha.23...v0.1.0-alpha.24](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.23...v0.1.0-alpha.24)

### Features

* **api:** manual updates ([1f158e2](https://github.com/stainless-api/stainless-api-cli/commit/1f158e21673453f97f17ca4935f5a1dad3f7ec0f))


### Chores

* bump version ([38b149b](https://github.com/stainless-api/stainless-api-cli/commit/38b149bd2a66c6604410c0924f710caa9866fb81))
* **internal:** version bump ([aef7105](https://github.com/stainless-api/stainless-api-cli/commit/aef7105c94e71f8c56c770be02c3b5f010904389))
* **internal:** version bump ([8d128be](https://github.com/stainless-api/stainless-api-cli/commit/8d128be1d6ef72eb019d80d165427f7f34d5e8d1))

## 0.1.0-alpha.23 (2025-06-30)

Full Changelog: [v0.1.0-alpha.22...v0.1.0-alpha.23](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.22...v0.1.0-alpha.23)

### Features

* improve project creation with config download and OpenAPI validation ([70f48ec](https://github.com/stainless-api/stainless-api-cli/commit/70f48ec8b2376826205c0d3fc631e819d5561ea7))


### Chores

* bump go sdk version ([fd6ea6a](https://github.com/stainless-api/stainless-api-cli/commit/fd6ea6a0d98f6e10232a44fa021c3cc5e0a77369))
* **ci:** only run for pushes and fork pull requests ([bd18456](https://github.com/stainless-api/stainless-api-cli/commit/bd184563ba5daade90637ccf945436f015e58d69))

## 0.1.0-alpha.22 (2025-06-26)

Full Changelog: [v0.1.0-alpha.21...v0.1.0-alpha.22](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.21...v0.1.0-alpha.22)

### Features

* add mcp command ([82f9ac7](https://github.com/stainless-api/stainless-api-cli/commit/82f9ac79d3025e10cab1f91f790d72ed75b30cb1))
* **api:** manual updates ([8809e1c](https://github.com/stainless-api/stainless-api-cli/commit/8809e1c1b087b2f87e57f20d341b0f19786e1d31))


### Bug Fixes

* change stainlessv0 import to stainless ([bb69895](https://github.com/stainless-api/stainless-api-cli/commit/bb698955e7a0107b45d4eef6d9fcdbcf820899dc))


### Chores

* update go sdk ([18af766](https://github.com/stainless-api/stainless-api-cli/commit/18af766bf7046273ff3e4d7ae0194962d7b6a8e2))

## 0.1.0-alpha.21 (2025-06-23)

Full Changelog: [v0.1.0-alpha.20...v0.1.0-alpha.21](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.20...v0.1.0-alpha.21)

### Features

* **api:** add diagnostics endpoint ([5bb2beb](https://github.com/stainless-api/stainless-api-cli/commit/5bb2beb07061e9699818db4858c1634223bcaf3d))


### Chores

* bump go sdk version ([a72a797](https://github.com/stainless-api/stainless-api-cli/commit/a72a79734aa7afd8fe25ae80ff963e688a55ede0))

## 0.1.0-alpha.20 (2025-06-23)

Full Changelog: [v0.1.0-alpha.19...v0.1.0-alpha.20](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.19...v0.1.0-alpha.20)

### Chores

* fix installation command in README ([b975888](https://github.com/stainless-api/stainless-api-cli/commit/b975888337c812bac92ad2fde27d4b5fa61a583f))
* remove openapi spec ([f4b0c0b](https://github.com/stainless-api/stainless-api-cli/commit/f4b0c0b892990d82f432c5c43ce7459bdda3febd))

## 0.1.0-alpha.19 (2025-06-20)

Full Changelog: [v0.1.0-alpha.18...v0.1.0-alpha.19](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.18...v0.1.0-alpha.19)

### Chores

* bump go package version ([abc91cb](https://github.com/stainless-api/stainless-api-cli/commit/abc91cb66f1e668c1e55f867ef2fbe4e22d6b61c))
* **internal:** codegen related update ([e84af71](https://github.com/stainless-api/stainless-api-cli/commit/e84af710aeaec0674afcef3d4eba4e5d3bbe5e71))

## 0.1.0-alpha.18 (2025-06-19)

Full Changelog: [v0.1.0-alpha.17...v0.1.0-alpha.18](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.17...v0.1.0-alpha.18)

### Chores

* run go mod tidy ([014b38b](https://github.com/stainless-api/stainless-api-cli/commit/014b38b3c5d50367a74eba773b8eb9c786325d5d))

## 0.1.0-alpha.17 (2025-06-19)

Full Changelog: [v0.1.0-alpha.16...v0.1.0-alpha.17](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.16...v0.1.0-alpha.17)

### Features

* **api:** make project nullable ([d51a660](https://github.com/stainless-api/stainless-api-cli/commit/d51a660892fa1b5196f37c8a4c297788c2ad7683))


### Chores

* bump go sdk version ([3865229](https://github.com/stainless-api/stainless-api-cli/commit/38652297409fbc08d4990fd631d0ee9d828f0005))
* **internal:** codegen related update ([5c194ff](https://github.com/stainless-api/stainless-api-cli/commit/5c194ff980c215b29ff16acef55cc0655d710537))

## 0.1.0-alpha.16 (2025-06-19)

Full Changelog: [v0.1.0-alpha.15...v0.1.0-alpha.16](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.15...v0.1.0-alpha.16)

### Features

* polish around logging ([38fec48](https://github.com/stainless-api/stainless-api-cli/commit/38fec484c830ae09ee32fd3f26ef19e6cf020eca))

## 0.1.0-alpha.15 (2025-06-18)

Full Changelog: [v0.1.0-alpha.14...v0.1.0-alpha.15](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.14...v0.1.0-alpha.15)

### Features

* add workspace status command ([adbc44f](https://github.com/stainless-api/stainless-api-cli/commit/adbc44f0539114437690008d87f16fa4a3ce0a4d))
* make branch parameter required for builds create ([0735363](https://github.com/stainless-api/stainless-api-cli/commit/07353638344ae98aaeff14e107730945ea3ff6ab))
* polish around logging ([e969a2d](https://github.com/stainless-api/stainless-api-cli/commit/e969a2dc7356dd4b9b0864f963a582dfbed43344))
* polish around logging ([a3fbc71](https://github.com/stainless-api/stainless-api-cli/commit/a3fbc71f38a195885b405aad8bec0c423694d692))


### Chores

* **internal:** codegen related update ([2dd834f](https://github.com/stainless-api/stainless-api-cli/commit/2dd834fed74e180a2ea70a4b8edf0bbb84a76d25))

## 0.1.0-alpha.14 (2025-06-17)

Full Changelog: [v0.1.0-alpha.13...v0.1.0-alpha.14](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.13...v0.1.0-alpha.14)

### Features

* add --openapi-spec and --stainless-config flags to workspace init ([530581e](https://github.com/stainless-api/stainless-api-cli/commit/530581e8ee3bd6d52b4a3745d229de40f87e4167))
* also automatically find openapi.json ([14eeeb4](https://github.com/stainless-api/stainless-api-cli/commit/14eeeb4c3635645059243d80509e4da7e05db28a))
* **api:** manual updates ([f4e6172](https://github.com/stainless-api/stainless-api-cli/commit/f4e61722f341bb7248e40dc22473408e7eef5b05))
* flesh out project create form ([c773199](https://github.com/stainless-api/stainless-api-cli/commit/c77319958ee6b285bd74191e1535a91bc77fbfce))
* sdkjson generation API ([9f908f1](https://github.com/stainless-api/stainless-api-cli/commit/9f908f1ff8be1c8455c533f0e494a5b6cc5cf4b2))

## 0.1.0-alpha.13 (2025-06-17)

Full Changelog: [v0.1.0-alpha.12...v0.1.0-alpha.13](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.12...v0.1.0-alpha.13)

### Chores

* **ci:** enable for pull requests ([3539b82](https://github.com/stainless-api/stainless-api-cli/commit/3539b82e2606cc2252ed808bf4feff883cab21c8))
* **internal:** codegen related update ([567a6e6](https://github.com/stainless-api/stainless-api-cli/commit/567a6e653a6ca70db50cf9cc80ac6deaa64be97a))

## 0.1.0-alpha.12 (2025-06-16)

Full Changelog: [v0.1.0-alpha.11...v0.1.0-alpha.12](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.11...v0.1.0-alpha.12)

### Features

* add build/upload steps to builds api ([4573934](https://github.com/stainless-api/stainless-api-cli/commit/45739343a28e9615e62e0c5126d1d66715212df9))
* add platform headers ([ca381b4](https://github.com/stainless-api/stainless-api-cli/commit/ca381b436607c3a7801b74e91a9147a2b59378bc))
* **api:** add v0 project create api ([a40eca6](https://github.com/stainless-api/stainless-api-cli/commit/a40eca6419d7ab8159491e861cde86d3e779b0a0))
* **api:** manual updates ([ef33c30](https://github.com/stainless-api/stainless-api-cli/commit/ef33c30512ca63d63380eb3f7c118e315696b2bc))
* **api:** manual updates ([80a6fb5](https://github.com/stainless-api/stainless-api-cli/commit/80a6fb50820cd4766359bcf74daccf72f44fad59))


### Bug Fixes

* changes har request format for snippets API some more ([4b2a898](https://github.com/stainless-api/stainless-api-cli/commit/4b2a898f947b7a4e000e758eb28cea5c00950687))
* fix type errors ([d5e1ae3](https://github.com/stainless-api/stainless-api-cli/commit/d5e1ae3f6c2566a42400a0a18394266a700f7463))


### Chores

* bump go package to 0.6.0 ([36a231d](https://github.com/stainless-api/stainless-api-cli/commit/36a231d0c7cb28e4111e9cbc3ac3c78833a340be))
* **internal:** codegen related update ([a108d1a](https://github.com/stainless-api/stainless-api-cli/commit/a108d1ad99e6bfd3469a8cecf48b4b299b199796))
* **internal:** codegen related update ([894c558](https://github.com/stainless-api/stainless-api-cli/commit/894c558dac1b497c69e1696453b9672bb923a0d9))


### Refactors

* move build_target_outputs to builds.target_outputs ([a085509](https://github.com/stainless-api/stainless-api-cli/commit/a08550995142a7f45786b9e33b3d20360068f2ec))

## 0.1.0-alpha.11 (2025-06-02)

Full Changelog: [v0.1.0-alpha.10...v0.1.0-alpha.11](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.10...v0.1.0-alpha.11)

### Features

* **api:** add diagnostics to build object ([fbbf7eb](https://github.com/stainless-api/stainless-api-cli/commit/fbbf7eba1eb3c06a9040a5e30904521a26021fa2))


### Chores

* make tap installation command shorter ([e8b060c](https://github.com/stainless-api/stainless-api-cli/commit/e8b060c69eb65b9828c863f97bd8fd4d064f7c0f))

## 0.1.0-alpha.10 (2025-05-30)

Full Changelog: [v0.1.0-alpha.9...v0.1.0-alpha.10](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.9...v0.1.0-alpha.10)

### Bug Fixes

* link auth and workspace subcommands to main CLI ([e73a657](https://github.com/stainless-api/stainless-api-cli/commit/e73a657b89ced3d651a82185755aa8d56c44c319))


### Chores

* add completions to .gitignore ([acf78a9](https://github.com/stainless-api/stainless-api-cli/commit/acf78a9bdd09c08d7c3846a465f50a7b1d070922))
* improve readme ([bf63975](https://github.com/stainless-api/stainless-api-cli/commit/bf63975c1fd32b59ae75bf1e6d345be584aae53d))

## 0.1.0-alpha.9 (2025-05-30)

Full Changelog: [v0.1.0-alpha.8...v0.1.0-alpha.9](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.8...v0.1.0-alpha.9)

### Features

* change --project-name flag to --project in workspace init command ([b1131e6](https://github.com/stainless-api/stainless-api-cli/commit/b1131e6259052b592af5b5a65f859357215a74e4))
* enhance workspace init with auto-detection and relative path resolution ([df091f0](https://github.com/stainless-api/stainless-api-cli/commit/df091f0fc109245407619bcae02e107ef5d6b28e))


### Chores

* add completions to release ([553394c](https://github.com/stainless-api/stainless-api-cli/commit/553394cbc05009671d5eef31e750a7b605ae7393))
* **internal:** codegen related update ([b90b5a6](https://github.com/stainless-api/stainless-api-cli/commit/b90b5a63cb1ecab77ac27ab3200b7cfa440e83c8))

## 0.1.0-alpha.8 (2025-05-29)

Full Changelog: [v0.1.0-alpha.7...v0.1.0-alpha.8](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.7...v0.1.0-alpha.8)

### Features

* **api:** manual updates ([27eb403](https://github.com/stainless-api/stainless-api-cli/commit/27eb403a9a2798555385c5ccb68c4b39080d40c0))


### Refactors

* move files into pkg/cmd ([1cdcbf4](https://github.com/stainless-api/stainless-api-cli/commit/1cdcbf46ca7f84c37b7e387a20f7a9c303498083))

## 0.1.0-alpha.7 (2025-05-28)

Full Changelog: [v0.1.0-alpha.6...v0.1.0-alpha.7](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.6...v0.1.0-alpha.7)

### Chores

* **internal:** codegen related update ([d50536d](https://github.com/stainless-api/stainless-api-cli/commit/d50536d13567bd394438d2521f00adb41d565037))

## 0.1.0-alpha.6 (2025-05-28)

Full Changelog: [v0.1.0-alpha.5...v0.1.0-alpha.6](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.5...v0.1.0-alpha.6)

### Features

* change projectName to project ([6ff942e](https://github.com/stainless-api/stainless-api-cli/commit/6ff942e73aa932b0986f13bdc28a6eef5f52e124))


### Chores

* add macos secrets to github workflwo ([cfb8c55](https://github.com/stainless-api/stainless-api-cli/commit/cfb8c55e2b5a61ffb27dd565f6d30d1c3705a993))

## 0.1.0-alpha.5 (2025-05-23)

Full Changelog: [v0.1.0-alpha.4...v0.1.0-alpha.5](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.4...v0.1.0-alpha.5)

### Features

* add preloading of the stainless config and openapi path ([37c089e](https://github.com/stainless-api/stainless-api-cli/commit/37c089eaf6df1150fa708cbccec3098b766e8417))
* add workspace command ([ad64eb7](https://github.com/stainless-api/stainless-api-cli/commit/ad64eb739a5ac2202e0fc7c720327945db7d679a))

## 0.1.0-alpha.4 (2025-05-23)

Full Changelog: [v0.1.0-alpha.3...v0.1.0-alpha.4](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.3...v0.1.0-alpha.4)

### Features

* **api:** change v0 path param projectName -&gt; project ([f4d71d2](https://github.com/stainless-api/stainless-api-cli/commit/f4d71d2a7329db59b1a9fa6d341619ab0532db6b))
* **api:** manual updates ([ad21d1e](https://github.com/stainless-api/stainless-api-cli/commit/ad21d1e0ee263bda24f88823c31f998295c842a2))
* **api:** manual updates ([2539725](https://github.com/stainless-api/stainless-api-cli/commit/253972559414ee57ef40cfcfc989fd5ea4e4cc27))


### Refactors

* minor refactor of method bodies ([8a3745d](https://github.com/stainless-api/stainless-api-cli/commit/8a3745dd7415ac36eeba2ec8b54e695845ffc9b5))

## 0.1.0-alpha.3 (2025-05-22)

Full Changelog: [v0.1.0-alpha.2...v0.1.0-alpha.3](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.2...v0.1.0-alpha.3)

### Features

* add readme ([517510b](https://github.com/stainless-api/stainless-api-cli/commit/517510b292283dff91a6f6403a133200050a49b2))
* **api:** add build compare to v0 ([ce9328e](https://github.com/stainless-api/stainless-api-cli/commit/ce9328e38387835ac55410e401194dc63d4e0fd9))
* **api:** bump go sdk version ([8500bf3](https://github.com/stainless-api/stainless-api-cli/commit/8500bf36df104935d6f63e3fd6ef651f21a253df))
* **api:** enable macos publishing ([c49364b](https://github.com/stainless-api/stainless-api-cli/commit/c49364b09ca3eb23f01ab254a1ec2c5b10a67a26))


### Bug Fixes

* fix custom code ([7d5ba04](https://github.com/stainless-api/stainless-api-cli/commit/7d5ba043c02683f4cf75c09029d65c136e39692e))
* update request schema for Postman webhook ([a9c5de3](https://github.com/stainless-api/stainless-api-cli/commit/a9c5de30d3ac155e9be266f29a35de13e7a84edc))


### Chores

* **internal:** codegen related update ([81448ca](https://github.com/stainless-api/stainless-api-cli/commit/81448cab092f0f1b50e1613c100f2fc2c97770d3))
* **internal:** codegen related update ([ac9557a](https://github.com/stainless-api/stainless-api-cli/commit/ac9557a75a32c0f838ab920648ef475969784cd8))
* **internal:** codegen related update ([4a1d74a](https://github.com/stainless-api/stainless-api-cli/commit/4a1d74ae1856a1375dbf74699d80823e5b5fb0cd))

## 0.1.0-alpha.2 (2025-04-30)

Full Changelog: [v0.1.0-alpha.1...v0.1.0-alpha.2](https://github.com/stainless-api/stainless-api-cli/compare/v0.1.0-alpha.1...v0.1.0-alpha.2)

### Features

* Add wait flag and polling mechanism for build completion ([960b063](https://github.com/stainless-api/stainless-api-cli/commit/960b063ce27de75fd4c847e69b78b2c0610c2ac1))


### Chores

* bump go sdk version ([3cbc3e4](https://github.com/stainless-api/stainless-api-cli/commit/3cbc3e46f7b8910037fa947e712e92bf05df5fc7))

## 0.1.0-alpha.1 (2025-04-30)

Full Changelog: [v0.0.1-alpha.0...v0.1.0-alpha.1](https://github.com/stainless-api/stainless-api-cli/compare/v0.0.1-alpha.0...v0.1.0-alpha.1)

### Features

* add typescript ([7a35228](https://github.com/stainless-api/stainless-api-cli/commit/7a3522871cb71d0511e8066d8c7052731d0fe56a))
* **api:** configs ([24361de](https://github.com/stainless-api/stainless-api-cli/commit/24361de65a1e8f72b1f7fda391e158f6737c7691))
* **api:** fix enum name conflict maybe ([e70ad16](https://github.com/stainless-api/stainless-api-cli/commit/e70ad1613201ddcbc500c372f904fd5dd6b33670))
* **api:** manual updates ([640c152](https://github.com/stainless-api/stainless-api-cli/commit/640c152a0fa8a2de1d8f90e6fe7588a73f8dd06b))
* **api:** manual updates ([d7fa9ae](https://github.com/stainless-api/stainless-api-cli/commit/d7fa9ae3231e881206843cd2044c1657f48c3455))
* **api:** manual updates ([aeaafee](https://github.com/stainless-api/stainless-api-cli/commit/aeaafee818795dbb22170c3eb778aa6a88ec3815))
* **api:** manual updates ([c694b82](https://github.com/stainless-api/stainless-api-cli/commit/c694b821b6720263e512c35d9e6c2fe71cb504f5))
* **api:** manual updates ([79ceae1](https://github.com/stainless-api/stainless-api-cli/commit/79ceae11587809ff66d89987cb43767890976727))
* **api:** manual updates ([c654f3d](https://github.com/stainless-api/stainless-api-cli/commit/c654f3dd0ebc143f724d7456920115e7eeff3543))
* **api:** manual updates ([69fd666](https://github.com/stainless-api/stainless-api-cli/commit/69fd666912c10edf82d9fac0b9a34e8a8201c0fd))
* **api:** manual updates ([514327b](https://github.com/stainless-api/stainless-api-cli/commit/514327b63421852a37eb64177e57987f64580620))
* **api:** manual updates ([1cc704d](https://github.com/stainless-api/stainless-api-cli/commit/1cc704d81c9bde5de6d4a0049d58868701919ec6))
* **api:** manual updates ([dfc8fa2](https://github.com/stainless-api/stainless-api-cli/commit/dfc8fa2e632544cf63d50ec02143ffe20cf30773))
* **api:** manual updates ([611beb0](https://github.com/stainless-api/stainless-api-cli/commit/611beb0f2fcff7fc46609864b0b51a102a825aad))
* **api:** manual updates ([0aeaa78](https://github.com/stainless-api/stainless-api-cli/commit/0aeaa78a2e8f59e8466cd0df451ee551245b4ef4))
* **api:** manual updates ([bc90b96](https://github.com/stainless-api/stainless-api-cli/commit/bc90b965b2df3c688bad5cef93cda86aaf228edf))
* **api:** manual updates ([15bda00](https://github.com/stainless-api/stainless-api-cli/commit/15bda00e341f32b2fa2a3ac059c311c014a17863))
* **api:** manual updates ([f598f10](https://github.com/stainless-api/stainless-api-cli/commit/f598f10545abbedca7db1a8e59c19a21f3393c44))
* **api:** manual updates ([947c01a](https://github.com/stainless-api/stainless-api-cli/commit/947c01ac614c6a0694cad2980150e9aaa6b8b76d))
* **api:** manual updates ([676ed23](https://github.com/stainless-api/stainless-api-cli/commit/676ed23f267362057b36797fed487db9beea3e26))
* **api:** parent build id ([540a1a0](https://github.com/stainless-api/stainless-api-cli/commit/540a1a06aef7294860f53326accc85379695b092))
* **api:** remove discriminator thing ([b87d671](https://github.com/stainless-api/stainless-api-cli/commit/b87d671eb9cb433d2e25e1e642bd65101f5c5012))
* **api:** rename api key ([b7e9c79](https://github.com/stainless-api/stainless-api-cli/commit/b7e9c790166bea892b438d58fe1c09a3fc9158ac))
* **api:** update via SDK Studio ([57b41f3](https://github.com/stainless-api/stainless-api-cli/commit/57b41f3587fb39b81719ec363c9c8a559c45eb70))
* **api:** update via SDK Studio ([8e817cc](https://github.com/stainless-api/stainless-api-cli/commit/8e817cc694231314fcd5f959497eeaade9819df5))
* **api:** update via SDK Studio ([8e4f061](https://github.com/stainless-api/stainless-api-cli/commit/8e4f061fa25a4d1847fc0025bc9dfb28fdc12f8a))
* **api:** update via SDK Studio ([#2](https://github.com/stainless-api/stainless-api-cli/issues/2)) ([10296a9](https://github.com/stainless-api/stainless-api-cli/commit/10296a982896077aa6923a67e697ba984177ef79))
* **api:** update via SDK Studio ([#4](https://github.com/stainless-api/stainless-api-cli/issues/4)) ([0380cc7](https://github.com/stainless-api/stainless-api-cli/commit/0380cc735be7df8284ab846dc1e78ab7d56f7942))
* **api:** use correct hashes ([fe39102](https://github.com/stainless-api/stainless-api-cli/commit/fe391029e2b2cfd7f56fbdeb7c6e3fe4b5efa86c))
* change list endpoint ([9943696](https://github.com/stainless-api/stainless-api-cli/commit/99436969643e83cf65d2cb9d06bf696225d0bc01))
* use urfave/cli library ([b1ae8e7](https://github.com/stainless-api/stainless-api-cli/commit/b1ae8e7ede73da7327e4b6764503e9152175923e))


### Bug Fixes

* don't overwrite headers ([#5](https://github.com/stainless-api/stainless-api-cli/issues/5)) ([2680134](https://github.com/stainless-api/stainless-api-cli/commit/2680134b981409ec9bdb5acaed2cf57c7cfb61be))


### Chores

* configure releases ([af7fad0](https://github.com/stainless-api/stainless-api-cli/commit/af7fad0fe84224aac40b03196412c3e285c57769))
* go live ([#1](https://github.com/stainless-api/stainless-api-cli/issues/1)) ([c9feb17](https://github.com/stainless-api/stainless-api-cli/commit/c9feb179487869081f1632c06f437220e74c5d5d))
* **internal:** codegen related update ([ec4f098](https://github.com/stainless-api/stainless-api-cli/commit/ec4f098eb4bb7a177acf30900dabc3b38ce08ca4))
* **internal:** codegen related update ([5d4d480](https://github.com/stainless-api/stainless-api-cli/commit/5d4d480936c1921673e32f21e92bcfaa63944913))


### Refactors

* split up completion target into multiple lines ([#3](https://github.com/stainless-api/stainless-api-cli/issues/3)) ([caa5a30](https://github.com/stainless-api/stainless-api-cli/commit/caa5a300c058762fababf68a5082c43c53f65b32))
