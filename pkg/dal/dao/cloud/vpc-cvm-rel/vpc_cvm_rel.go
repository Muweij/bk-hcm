/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 混合云管理平台 (BlueKing - Hybrid Cloud Management System) available.
 * Copyright (C) 2022 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 *
 * to the current version of the project delivered to anyone in the future.
 */

package vpccvmrel

import (
	"fmt"

	"hcm/pkg/api/core"
	"hcm/pkg/criteria/errf"
	vpc "hcm/pkg/dal/dao/cloud"
	"hcm/pkg/dal/dao/cloud/cvm"
	"hcm/pkg/dal/dao/orm"
	"hcm/pkg/dal/dao/tools"
	"hcm/pkg/dal/dao/types"
	"hcm/pkg/dal/table"
	"hcm/pkg/dal/table/cloud"
	"hcm/pkg/kit"
	"hcm/pkg/logs"
	"hcm/pkg/runtime/filter"

	"github.com/jmoiron/sqlx"
)

// Interface only used for vpc and cvm relation.
type Interface interface {
	BatchCreateWithTx(kt *kit.Kit, tx *sqlx.Tx, rels []cloud.VpcCvmRelTable) error
	List(kt *kit.Kit, opt *types.ListOption) (*types.ListVpcCvmRelDetails, error)
	DeleteWithTx(kt *kit.Kit, tx *sqlx.Tx, expr *filter.Expression) error
	ListJoinVpc(kt *kit.Kit, cvmIDs []string) (*types.ListVpcCvmRelsJoinVpcDetails, error)
}

var _ Interface = new(Dao)

// Dao define vpc and cvm relation dao.
type Dao struct {
	Orm orm.Interface
}

// ListJoinVpc list vpc cvm rel with vpc detail.
func (dao Dao) ListJoinVpc(kt *kit.Kit, cvmIDs []string) (*types.ListVpcCvmRelsJoinVpcDetails, error) {
	if len(cvmIDs) == 0 {
		return nil, errf.Newf(errf.InvalidParameter, "cvm ids is required")
	}

	sql := fmt.Sprintf(`SELECT %s, %s FROM %s as rel left join %s as vpc on rel.vpc_id = vpc.id 
        where cvm_id in (:cvm_ids)`,
		cloud.VpcColumns.FieldsNamedExprWithout(types.DefaultRelJoinWithoutField),
		tools.BaseRelJoinSqlBuild("rel", "vpc", "vpc_id", "cvm_id"),
		table.VpcCvmRelTable, table.VpcTable)

	details := make([]types.VpcWithCvmID, 0)
	if err := dao.Orm.Do().Select(kt.Ctx, &details, sql, map[string]interface{}{"cvm_ids": cvmIDs}); err != nil {
		logs.ErrorJson("select vpc cvm rels join vpc failed, err: %v, sql: (%s), rid: %s", err, sql, kt.Rid)
		return nil, err
	}

	return &types.ListVpcCvmRelsJoinVpcDetails{Details: details}, nil
}

// BatchCreateWithTx batch create vpc cvm rel with transaction.
func (dao Dao) BatchCreateWithTx(kt *kit.Kit, tx *sqlx.Tx, rels []cloud.VpcCvmRelTable) error {
	// 校验关联资源是否存在
	vpcIDs := make([]string, 0)
	cvmIDs := make([]string, 0)
	for _, rel := range rels {
		vpcIDs = append(vpcIDs, rel.VpcID)
		cvmIDs = append(cvmIDs, rel.CvmID)
	}

	vpcMap, err := vpc.ListVpc(kt, dao.Orm, vpcIDs)
	if err != nil {
		logs.Errorf("list vpc failed, err: %v, ids: %v, rid: %s", err, vpcIDs, kt.Rid)
		return err
	}

	if len(vpcMap) != len(vpcIDs) {
		logs.Errorf("get vpc count not right, err: %v, ids: %v, count: %d, rid: %s", err, vpcIDs, len(vpcMap), kt.Rid)
		return fmt.Errorf("get vpc count not right")
	}

	cvmMap, err := cvm.ListCvm(kt, dao.Orm, cvmIDs)
	if err != nil {
		logs.Errorf("list vpc failed, err: %v, ids: %v, rid: %s", err, vpcIDs, kt.Rid)
		return err
	}

	if len(cvmMap) != len(cvmIDs) {
		logs.Errorf("get cvm count not right, err: %v, ids: %v, count: %d, rid: %s", err, cvmIDs, len(cvmMap), kt.Rid)
		return fmt.Errorf("get cvm count not right")
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s)	VALUES(%s)`, table.VpcCvmRelTable, cloud.VpcCvmRelColumns.ColumnExpr(),
		cloud.VpcCvmRelColumns.ColonNameExpr())

	if err := dao.Orm.Txn(tx).BulkInsert(kt.Ctx, sql, rels); err != nil {
		logs.Errorf("insert %s failed, err: %v, rid: %s", table.VpcCvmRelTable, err, kt.Rid)
		return fmt.Errorf("insert %s failed, err: %v", table.VpcCvmRelTable, err)
	}

	return nil
}

// List vpc cvm rel.
func (dao Dao) List(kt *kit.Kit, opt *types.ListOption) (*types.ListVpcCvmRelDetails, error) {
	if opt == nil {
		return nil, errf.New(errf.InvalidParameter, "list options is nil")
	}

	if err := opt.Validate(filter.NewExprOption(filter.RuleFields(cloud.VpcCvmRelColumns.ColumnTypes())),
		core.DefaultPageOption); err != nil {
		return nil, err
	}

	whereExpr, whereValue, err := opt.Filter.SQLWhereExpr(tools.DefaultSqlWhereOption)
	if err != nil {
		return nil, err
	}

	if opt.Page.Count {
		sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, table.VpcCvmRelTable, whereExpr)

		count, err := dao.Orm.Do().Count(kt.Ctx, sql, whereValue)
		if err != nil {
			logs.ErrorJson("count vpc cvm rels failed, err: %v, filter: %s, rid: %s", err,
				opt.Filter, kt.Rid)
			return nil, err
		}

		return &types.ListVpcCvmRelDetails{Count: count}, nil
	}

	pageExpr, err := types.PageSQLExpr(opt.Page, types.DefaultPageSQLOption)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s %s %s`, cloud.VpcCvmRelColumns.FieldsNamedExpr(opt.Fields),
		table.VpcCvmRelTable, whereExpr, pageExpr)

	details := make([]cloud.VpcCvmRelTable, 0)
	if err = dao.Orm.Do().Select(kt.Ctx, &details, sql, whereValue); err != nil {
		logs.ErrorJson("select vpc cvm rels failed, err: %v, filter: %s, rid: %s", err, opt.Filter, kt.Rid)
		return nil, err
	}

	return &types.ListVpcCvmRelDetails{Details: details}, nil
}

// DeleteWithTx delete vpc cvm rel with transaction.
func (dao Dao) DeleteWithTx(kt *kit.Kit, tx *sqlx.Tx, expr *filter.Expression) error {
	if expr == nil {
		return errf.New(errf.InvalidParameter, "filter expr is required")
	}

	whereExpr, whereValue, err := expr.SQLWhereExpr(tools.DefaultSqlWhereOption)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`DELETE FROM %s %s`, table.VpcCvmRelTable, whereExpr)
	if _, err = dao.Orm.Txn(tx).Delete(kt.Ctx, sql, whereValue); err != nil {
		logs.ErrorJson("delete vpc cvm rels failed, err: %v, filter: %s, rid: %s", err, expr, kt.Rid)
		return err
	}

	return nil
}
