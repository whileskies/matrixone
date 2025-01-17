# Export data in general method

This document will guide you export data in general method.

## Before you start

- Make sure you have already [installed and launched MatrixOne using source](https://docs.matrixorigin.io/0.5.1/MatrixOne/Get-Started/install-standalone-matrixone/#method-1-building-from-source) or [installed and launched MatrixOne using binary packages](https://docs.matrixorigin.io/0.5.1/MatrixOne/Get-Started/install-standalone-matrixone/#method-2-downloading-binary-packages).

- Make sure you have already [connected MatrixOne Server](../../Get-Started/connect-to-matrixone-server.md).

!!! note
    MatrixOne does not support using `dump` to export the tables, only supports using `SELECT...INTO OUTFILE` to export the table.

- **Scenario description**：Create a table in MatrixOne, export the table to your local path, for example: *~/tmp/export_demo/export_datatable.txt*.

### Steps

1. Create tables in MatrixOne:

    ```sql
    create database aaa;
    use aaa;
    CREATE TABLE `user` (`id` int(11) ,`user_name` varchar(255) ,`sex` varchar(255));
    insert into user(id,user_name,sex) values('1', 'weder', 'man'), ('2', 'tom', 'man'), ('3', 'wederTom', 'man');
    select * from user;
    +------+-----------+------+
    | id   | user_name | sex  |
    +------+-----------+------+
    |    1 | weder     | man  |
    |    2 | tom       | man  |
    |    3 | wederTom  | man  |
    +------+-----------+------+
    ```

2. Export the table to your local directory, for example, *~/tmp/export_demo/export_datatable.txt*

    ```
    select * from user into outfile '~/tmp/export_demo/export_datatable.txt'
    ```

3. Check the table in your local directory *~/tmp/export_demo/export_datatable.txt*.

    Open the *~/tmp/export_demo/* directory, *export_datatable.txt* file is created, then open the file, check the result as below:

    ```
    id,user_name,sex
    1,"weder","man"
    2,"tom","man"
    3,"wederTom","man"
    ```
