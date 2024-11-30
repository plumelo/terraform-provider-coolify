package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
)

type postgresqlDatabaseModel struct {
	commonDatabaseModel
	PostgresConf           types.String `tfsdk:"postgres_conf"`
	PostgresDb             types.String `tfsdk:"postgres_db"`
	PostgresHostAuthMethod types.String `tfsdk:"postgres_host_auth_method"`
	PostgresInitdbArgs     types.String `tfsdk:"postgres_initdb_args"`
	PostgresPassword       types.String `tfsdk:"postgres_password"`
	PostgresUser           types.String `tfsdk:"postgres_user"`
}

func (m postgresqlDatabaseModel) FromAPI(apiModel *api.Database, state postgresqlDatabaseModel) (postgresqlDatabaseModel, error) {
	apiModel.ValueByDiscriminator()
	db, err := apiModel.AsPostgresqlDatabase()
	if err != nil {
		return postgresqlDatabaseModel{}, err
	}

	return postgresqlDatabaseModel{
		commonDatabaseModel:    commonDatabaseModel{}.FromAPI(apiModel, state.commonDatabaseModel),
		PostgresConf:           optionalString(db.PostgresConf),
		PostgresDb:             optionalString(db.PostgresDb),
		PostgresHostAuthMethod: optionalString(db.PostgresHostAuthMethod),
		PostgresInitdbArgs:     optionalString(db.PostgresInitdbArgs),
		PostgresPassword:       optionalString(db.PostgresPassword),
		PostgresUser:           optionalString(db.PostgresUser),
	}, nil
}
