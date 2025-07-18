// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: recharges.sql

package zing

import (
	"context"
	"database/sql"
)

const getRecharge = `-- name: GetRecharge :one
SELECT
	id, customerid, productid, method, paymentservicepaymentid, paymentservicepayload, paymentservicequeryparams, rechargesuccessful, failurereason, paymentamount, expirydate, previousrmexpirydate, userid, fromrmsvcid, tormsvcid, datecreated, deleted
FROM
	Recharges
WHERE
	id = ?
LIMIT
	1
`

func (q *Queries) GetRecharge(ctx context.Context, id string) (Recharge, error) {
	row := q.db.QueryRowContext(ctx, getRecharge, id)
	var i Recharge
	err := row.Scan(
		&i.ID,
		&i.Customerid,
		&i.Productid,
		&i.Method,
		&i.Paymentservicepaymentid,
		&i.Paymentservicepayload,
		&i.Paymentservicequeryparams,
		&i.Rechargesuccessful,
		&i.Failurereason,
		&i.Paymentamount,
		&i.Expirydate,
		&i.Previousrmexpirydate,
		&i.Userid,
		&i.Fromrmsvcid,
		&i.Tormsvcid,
		&i.Datecreated,
		&i.Deleted,
	)
	return i, err
}

const getRechargeTypeCounts = `-- name: GetRechargeTypeCounts :many
SELECT
	recharge_name, recharge_count, recharge_period, recharge_max_date
FROM
	(
		SELECT
			t3.Name AS recharge_name,
			COUNT(*) AS recharge_count,
			CASE
				WHEN ? = 'weeks' THEN CONCAT(
					FLOOR((DAY(t1.DateCreated) - 1) / 7) + 1,
					'-',
					MONTH(t1.DateCreated),
					'-',
					YEAR(t1.DateCreated)
				)
				WHEN ? = 'months' THEN CONCAT(MONTH(t1.DateCreated), '-', YEAR(t1.DateCreated))
			END AS recharge_period,
			MAX(t1.DateCreated) AS recharge_max_date
		FROM
			Recharges t1
			LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
			LEFT JOIN Products t3 ON t1.ProductId = t3.Id
            LEFT JOIN Addresses t4 ON t2.AddressId = t4.Id
		WHERE
			TRIM(LOWER(t4.POP)) LIKE CONCAT(TRIM(LOWER(?)), '%')
			AND(
                (
                    ? = 'weeks'
                    AND t1.DateCreated >= 
                        CASE 
                            WHEN ? = 1 THEN DATE_FORMAT(NOW(), '%Y-%m-01 00:00:00')
                            ELSE DATE_FORMAT(DATE_SUB(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL (? - 1) WEEK), '%Y-%m-01 00:00:00')
                        END
                )
                OR(
                    ? = 'months'
                    AND t1.DateCreated >= 
                        CASE 
                            WHEN ? = 1 THEN DATE_FORMAT(NOW(), '%Y-%m-01 00:00:00')
                            ELSE DATE_FORMAT(DATE_SUB(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL (? - 1) MONTH), '%Y-%m-01 00:00:00')
                        END
                )
			)
		GROUP BY
			recharge_name,
			recharge_period
	) AS sub
ORDER BY
	recharge_max_date ASC,
    recharge_count DESC
`

type GetRechargeTypeCountsParams struct {
	Period interface{}
	Poi    string
	Count  interface{}
}

type GetRechargeTypeCountsRow struct {
	RechargeName    sql.NullString
	RechargeCount   int64
	RechargePeriod  interface{}
	RechargeMaxDate interface{}
}

func (q *Queries) GetRechargeTypeCounts(ctx context.Context, arg GetRechargeTypeCountsParams) ([]GetRechargeTypeCountsRow, error) {
	rows, err := q.db.QueryContext(ctx, getRechargeTypeCounts,
		arg.Period,
		arg.Period,
		arg.Poi,
		arg.Period,
		arg.Count,
		arg.Count,
		arg.Period,
		arg.Count,
		arg.Count,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRechargeTypeCountsRow
	for rows.Next() {
		var i GetRechargeTypeCountsRow
		if err := rows.Scan(
			&i.RechargeName,
			&i.RechargeCount,
			&i.RechargePeriod,
			&i.RechargeMaxDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecharges = `-- name: GetRecharges :many
SELECT
	id, customerid, productid, method, paymentservicepaymentid, paymentservicepayload, paymentservicequeryparams, rechargesuccessful, failurereason, paymentamount, expirydate, previousrmexpirydate, userid, fromrmsvcid, tormsvcid, datecreated, deleted
FROM
	Recharges
LIMIT
	?
OFFSET
	?
`

type GetRechargesParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetRecharges(ctx context.Context, arg GetRechargesParams) ([]Recharge, error) {
	rows, err := q.db.QueryContext(ctx, getRecharges, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recharge
	for rows.Next() {
		var i Recharge
		if err := rows.Scan(
			&i.ID,
			&i.Customerid,
			&i.Productid,
			&i.Method,
			&i.Paymentservicepaymentid,
			&i.Paymentservicepayload,
			&i.Paymentservicequeryparams,
			&i.Rechargesuccessful,
			&i.Failurereason,
			&i.Paymentamount,
			&i.Expirydate,
			&i.Previousrmexpirydate,
			&i.Userid,
			&i.Fromrmsvcid,
			&i.Tormsvcid,
			&i.Datecreated,
			&i.Deleted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
