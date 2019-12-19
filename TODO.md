* add user agent that mentions terraform provider so hurricane can see it in logs

* support AAAA records (need placeholder address like for A records)

* include tests (right now these are manual)
    - create when there is no A record at all - should FAIL
    - create when A record exists with a non-placeholder address - should FAIL
    - create when the A record exists with placeholder address - should update

* "import" support for existing records
