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
