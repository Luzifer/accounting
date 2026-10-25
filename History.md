# 0.8.0 / 2026-10-25

* Improvements
  * chore: add develop-env and healthz endpoint
  * chore: centralize fetch logic
  * chore: extract account create modal
  * chore: extract money transfer modal
  * chore: migrate to modal-host component
  * chore: port frontend to Typescript
  * chore: separate components from views
  * chore: sort functions
  * chore: switch to scratch final stage
  * fix: allow jumping from the category with enter
  * fix: when creating a transfer treat current account as from

* Dependencies
  * chore(deps): update dependency @babel/eslint-parser to v7.28.6
  * chore(deps): update dependency esbuild-plugin-vue3 to v0.5.1
  * chore(deps): update dependency esbuild-sass-plugin to v3.7.0
  * chore(deps): update dependency esbuild to v0.28.2 (#141)
  * chore(deps): update dependency eslint-plugin-vue to v10.9.2 (#86)
  * chore(deps): update dependency eslint to v10.9.1 (#158)
  * chore(deps): update dependency globals to v17.9.0 (#136)
  * chore(deps): update dependency @types/bootstrap to v5.2.11 (#80)
  * chore(deps): update dependency typescript to v6.0.3
  * chore(deps): update dependency vue-eslint-parser to v10.4.1 (#83)
  * chore(deps): update ghcr.io/luzifer-docker/pnpm docker tag to v11.28.6 (#196)
  * chore(deps): update golang docker tag to v1.27.2 (#193)
  * chore(deps): update go toolchain directive to v1.27.3 [security] (#200)
  * chore(deps): update module golang.org/x/text to v0.42.0 [security] (#188)
  * chore(deps): update typescript-eslint monorepo to v8.72.0 (#201)
  * chore(sec): upgrade transitive dependencies
  * fix(deps): update dependency bootstrap to v5.3.8
  * fix(deps): update dependency @fortawesome/fontawesome-free to v7.3.1 (#118)
  * fix(deps): update dependency vue-router to v5.4.0 (#192)
  * fix(deps): update dependency vue to v3.5.44 (#195)
  * fix(deps): update module github.com/luzifer/go_helpers/* (#175)
  * fix(deps): update module github.com/luzifer/rconfig/v2 to v2.6.2
  * fix(deps): update module github.com/sirupsen/logrus to v1.10.2 (#161)
  * fix(deps): update module github.com/stretchr/testify to v1.12.1 (#153)
  * fix(deps): update module gopkg.in/evanphx/json-patch.v5 to v5.9.11
  * fix(deps): update module gorm.io/driver/postgres to v1.6.3 (#176)
  * fix(deps): update module gorm.io/gorm to v1.31.2 (#97)

# 0.7.1 / 2024-12-12

  * Update Go dependencies

# 0.7.0 / 2024-12-07

  * Improvements
    * Add category-activity transaction view

  * Bugfixes
    * Fix number formatting for 1000, 10000, ...
    * Fix: Auto-clear money transfers between categories
    * Update dependencies

# 0.6.0 / 2024-02-09

  * Add "clear today" functionality
  * Fix: Paired transactions must be updated with negative sum

# 0.5.0 / 2024-02-06

  * Add description to money transfers
  * Fix: Broken category ID parsing for money transfers

# 0.4.0 / 2024-02-03

  * Move creation of starting balances to backend
  * Restore selected date-range for transaction list
  * Fix: Update amount of paired transactions
  * Fix: Unallocated gets unreadable if negative

# 0.3.0 / 2024-02-01

  * Add account reconcilation
  * Sort uncleared transactions to the top
  * Fix: Allow updating transfer transactions
  * Fix: Clear editor after saving new transaction
  * Fix: Submit account ID when creating transaction
  * CI: Add test / build / publish

# 0.2.0 / 2024-01-28

  * Improvements
    * Define transfer-money transactions as pair
    * Improve display of zero values

# 0.1.0 / 2024-01-22

  * Initial release
