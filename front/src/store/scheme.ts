import http from '@/http';
import { defineStore } from 'pinia';
import { QueryFilterType, IPageQuery } from '@/typings/common';
import { IAreaInfo, IBizTypeResData, ICountriesListResData, IGenerateSchemesResData, IUserDistributionResData, IGenerateSchemesReqParams, IRecommendSchemeList } from '@/typings/scheme';

const { BK_HCM_AJAX_URL_PREFIX } = window.PROJECT_CONFIG;


// 资源选型模块相关状态管理和接口定义
export const useSchemeStore = defineStore({
  id: 'schemeStore',
  state: () => ({
    userDistribution: [] as Array<IAreaInfo>,
    recommendationSchemes: [] as IRecommendSchemeList,
  }),
  actions: {
    setUserDistribution(data: Array<IAreaInfo>) {
      this.userDistribution = data;
    },
    setRecommendationSchemes(data: IRecommendSchemeList) {
      this.recommendationSchemes = data;
    },
    /**
     * 获取资源选型方案列表
     * @param filter 过滤参数
     * @param page 分页参数
     * @returns
     */
    listCloudSelectionScheme (filter: QueryFilterType, page: IPageQuery) {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/schemes/list`, { filter, page });
    },
    /**
     * 删除资源选型方案
     * @param ids 方案id列表
     * @returns 
     */
    deleteCloudSelectionScheme (ids: string[]) {
      return http.delete(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/schemes/batch`, { data: { ids } });
    },
    /**
     * 获取资源选型方案详情
     * @param id 方案id
     * @returns 
     */
    getCloudSelectionScheme (id: string) {
      return http.get(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/schemes/${id}`);
    },
    /**
     * 更新资源选型方案
     * @param id 方案id
     * @param data 方案数据
     */
    updateCloudSelectionScheme (id: string, data: { name: string; bk_biz_id?: number; }) {
      return http.patch(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/schemes/${id}`, data);
    },
    /**
     * 获取收藏的资源选型方案列表
     * @returns
     */
    listCollection () {
      return http.get(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/collections/cloud_selection_scheme/list`);
    },
    /** 添加收藏
    * @param id 方案id
    * @returns
    */
    createCollection (id: string) {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/collections/create`, { res_type: 'cloud_selection_scheme', res_id: id });
    },
    /**
      * 取消收藏
      * @param id 收藏id
      * @returns
      */
    deleteCollection (id: number) {
      return http.delete(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/collections/${id}`);
    },
    /**
     * 查询IDC机房列表
     * @param filter 过滤参数
     * @param page 分页参数
     * @returns 
     */
    listIdc (filter: QueryFilterType, page: IPageQuery) {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/idcs/list`, { filter, page });
    },
    /**
     * 查询业务延迟数据
     * @param topo 拓扑列表
     * @param ids idc列表
     */
    queryBizLatency(topo: IAreaInfo[], ids: string[]) {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/latency/biz/query`, { area_topo: topo, idc_ids: ids });
    },
    /**
     * 查询ping延迟数据
     * @param topo 拓扑列表
     * @param ids idc列表
     */
    queryPingLatency(topo: IAreaInfo[], ids: string[]) {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/latency/ping/query`, { area_topo: topo, idc_ids: ids });
    },
    /**
     * 获取云选型数据支持的国家列表
     * @returns 
     */
    listCountries (): ICountriesListResData {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/countries/list`);
    },
    /**
     * 获取业务类型列表
     * @param page 分页参数
     * @returns 
     */
    listBizTypes (page: IPageQuery): IBizTypeResData {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/biz_types/list`, { page });
    },
    /**
     * 获取云选型用户分布占比
     * @param area_topo 需要查询的国家列表
     * @returns 
     */
    queryUserDistributions (area_topo: Array<IAreaInfo>): IUserDistributionResData {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/user_distributions/query`, { area_topo });
    },
    /**
     * 生成云资源选型方案
     * @param formData 业务属性
     * @returns 
     */
    generateSchemes (data: IGenerateSchemesReqParams): IGenerateSchemesResData {
      return http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/selections/schemes/generate`, data);
    }
  },
});
