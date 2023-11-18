package system

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateCasbin
//@description: 更新casbin权限
//@param: authorityId string, casbinInfos []request.CasbinInfo
//@return: error

type CasbinService struct{}

var CasbinServiceApp = new(CasbinService)

func (casbinService *CasbinService) UpdateCasbin(AuthorityID uint, casbinInfos []request.CasbinInfo) error {
	authorityId := strconv.Itoa(int(AuthorityID))
	casbinService.ClearCasbin(0, authorityId)
	rules := [][]string{}
	//做权限去重处理
	deduplicateMap := make(map[string]bool)
	for _, v := range casbinInfos {
		key := authorityId + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
	}
	e := casbinService.Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateCasbinApi
//@description: API更新随动
//@param: oldPath string, newPath string, oldMethod string, newMethod string
//@return: error

func (casbinService *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := global.GVA_DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	e := casbinService.Casbin()
	err = e.LoadPolicy()
	if err != nil {
		return err
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetPolicyPathByAuthorityId
//@description: 获取权限列表
//@param: authorityId string
//@return: pathMaps []request.CasbinInfo

func (casbinService *CasbinService) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []request.CasbinInfo) {
	e := casbinService.Casbin()
	authorityId := strconv.Itoa(int(AuthorityID))
	list := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, request.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ClearCasbin
//@description: 清除匹配的权限
//@param: v int, p ...string
//@return: bool

func (casbinService *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := casbinService.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Casbin
//@description: 持久化到数据库  引入自定义规则
//@return: *casbin.Enforcer

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

func (casbinService *CasbinService) Casbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		a, err := gormadapter.NewAdapterByDB(global.GVA_DB)
		if err != nil {
			zap.L().Error("适配数据库失败请检查casbin表是否为InnoDB引擎!", zap.Error(err))
			return
		}
		text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		m, err := model.NewModelFromString(text)
		if err != nil {
			zap.L().Error("字符串加载模型失败!", zap.Error(err))
			return
		}
		syncedCachedEnforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
		syncedCachedEnforcer.SetExpireTime(60 * 60)
		_ = syncedCachedEnforcer.LoadPolicy()
	})
	return syncedCachedEnforcer
}

func (casbinService *CasbinService) CasbinWithDB(db *gorm.DB) (*casbin.SyncedEnforcer, error) {
	var a *gormadapter.Adapter
	var m model.Model
	var s *casbin.SyncedEnforcer
	var err error
	var text = `
			[request_definition]
			r = sub, obj, act
			
			[policy_definition]
			p = sub, obj, act
			
			[role_definition]
			g = _, _
			
			[policy_effect]
			e = some(where (p.eft == allow))
			
			[matchers]
			m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
			`

	if a, err = gormadapter.NewAdapterByDB(db); err != nil {
		zap.L().Error("适配数据库失败请检查casbin表是否为InnoDB引擎!", zap.Error(err))
		return nil, err
	}

	if m, err = model.NewModelFromString(text); err != nil {
		zap.L().Error("字符串加载模型失败!", zap.Error(err))
		return nil, err
	}

	if s, err = casbin.NewSyncedEnforcer(m, a); err != nil {
		zap.L().Error("通过文件或 DB 创建执行器失败!", zap.Error(err))
		return nil, err
	}

	if err = s.LoadPolicy(); err != nil {
		zap.L().Error("从持久层中加载policy规则失败!", zap.Error(err))
		return nil, err
	}

	return s, nil
}
