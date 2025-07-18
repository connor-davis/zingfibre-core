// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: rm_invoices.sql

package radius

import (
	"context"
)

const getRmInvoice = `-- name: GetRmInvoice :one
SELECT
    id, invgroup, invnum, managername, username, date, bytesdl, bytesul, bytescomb, downlimit, uplimit, comblimit, time, uptimelimit, days, expiration, capdl, capul, captotal, captime, capdate, service, comment, transid, amount, address, city, zip, country, state, fullname, taxid, contractid, paymentopt, invtype, paymode, paid, price, tax, advtax, vatpercent, advtaxpercent, remark, balance, gwtransid, phone, mobile
FROM
    rm_invoices
WHERE
    id = ?
LIMIT
    1
`

func (q *Queries) GetRmInvoice(ctx context.Context, id int32) (RmInvoice, error) {
	row := q.db.QueryRowContext(ctx, getRmInvoice, id)
	var i RmInvoice
	err := row.Scan(
		&i.ID,
		&i.Invgroup,
		&i.Invnum,
		&i.Managername,
		&i.Username,
		&i.Date,
		&i.Bytesdl,
		&i.Bytesul,
		&i.Bytescomb,
		&i.Downlimit,
		&i.Uplimit,
		&i.Comblimit,
		&i.Time,
		&i.Uptimelimit,
		&i.Days,
		&i.Expiration,
		&i.Capdl,
		&i.Capul,
		&i.Captotal,
		&i.Captime,
		&i.Capdate,
		&i.Service,
		&i.Comment,
		&i.Transid,
		&i.Amount,
		&i.Address,
		&i.City,
		&i.Zip,
		&i.Country,
		&i.State,
		&i.Fullname,
		&i.Taxid,
		&i.Contractid,
		&i.Paymentopt,
		&i.Invtype,
		&i.Paymode,
		&i.Paid,
		&i.Price,
		&i.Tax,
		&i.Advtax,
		&i.Vatpercent,
		&i.Advtaxpercent,
		&i.Remark,
		&i.Balance,
		&i.Gwtransid,
		&i.Phone,
		&i.Mobile,
	)
	return i, err
}

const getRmInvoices = `-- name: GetRmInvoices :many
SELECT
    id, invgroup, invnum, managername, username, date, bytesdl, bytesul, bytescomb, downlimit, uplimit, comblimit, time, uptimelimit, days, expiration, capdl, capul, captotal, captime, capdate, service, comment, transid, amount, address, city, zip, country, state, fullname, taxid, contractid, paymentopt, invtype, paymode, paid, price, tax, advtax, vatpercent, advtaxpercent, remark, balance, gwtransid, phone, mobile
FROM
    rm_invoices
LIMIT
    ?
OFFSET
    ?
`

type GetRmInvoicesParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetRmInvoices(ctx context.Context, arg GetRmInvoicesParams) ([]RmInvoice, error) {
	rows, err := q.db.QueryContext(ctx, getRmInvoices, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RmInvoice
	for rows.Next() {
		var i RmInvoice
		if err := rows.Scan(
			&i.ID,
			&i.Invgroup,
			&i.Invnum,
			&i.Managername,
			&i.Username,
			&i.Date,
			&i.Bytesdl,
			&i.Bytesul,
			&i.Bytescomb,
			&i.Downlimit,
			&i.Uplimit,
			&i.Comblimit,
			&i.Time,
			&i.Uptimelimit,
			&i.Days,
			&i.Expiration,
			&i.Capdl,
			&i.Capul,
			&i.Captotal,
			&i.Captime,
			&i.Capdate,
			&i.Service,
			&i.Comment,
			&i.Transid,
			&i.Amount,
			&i.Address,
			&i.City,
			&i.Zip,
			&i.Country,
			&i.State,
			&i.Fullname,
			&i.Taxid,
			&i.Contractid,
			&i.Paymentopt,
			&i.Invtype,
			&i.Paymode,
			&i.Paid,
			&i.Price,
			&i.Tax,
			&i.Advtax,
			&i.Vatpercent,
			&i.Advtaxpercent,
			&i.Remark,
			&i.Balance,
			&i.Gwtransid,
			&i.Phone,
			&i.Mobile,
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
