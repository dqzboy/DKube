<template>
    <div class="daemonset">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="daemonset-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="6">
                                <div>
                                    <span>命名空间: </span>
                                    <el-select v-model="namespaceValue" filterable placeholder="请选择">
                                        <el-option
                                        v-for="(item, index) in namespaceList"
                                        :key="index"
                                        :label="item.metadata.name"
                                        :value="item.metadata.name">
                                        </el-option>
                                    </el-select>
                                </div>
                            </el-col>
                            <el-col :span="2" :offset="16">
                                <div>
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getDaemonSets()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="daemonset-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button disabled style="border-radius:2px;" icon="Edit" type="primary">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="daemonset-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getDaemonSets()">搜索</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="daemonset-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="daemonSetList"
                        v-loading="appLoading">
                            <el-table-column width="20"></el-table-column>
                            <el-table-column align=left label="DaemonSet名">
                                <template v-slot="scope">
                                    <a class="daemonset-body-daemonsetname">{{ scope.row.metadata.name }}</a>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="标签">
                                <template v-slot="scope">
                                    <div v-for="(val, key) in scope.row.metadata.labels" :key="key">
                                        <el-popover
                                            placement="right"
                                            :width="200"
                                            trigger="hover"
                                            :content="key + ':' + val">
                                            <template #reference>
                                                <el-tag style="margin-bottom: 5px" type="warning">{{ ellipsis(key + ":" + val) }}</el-tag>
                                            </template>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="容器组">
                                <template v-slot="scope">
                                    <span>{{ scope.row.status.numberAvailable>0?scope.row.status.numberAvailable:0  }} / {{ scope.row.status.desiredNumberScheduled>0?scope.row.status.desiredNumberScheduled:0 }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="镜像">
                                <template v-slot="scope">
                                    <div v-for="(val, key) in scope.row.spec.template.spec.containers" :key="key">
                                        <el-popover
                                            placement="right"
                                            :width="200"
                                            trigger="hover"
                                            :content="val.image">
                                            <template #reference>
                                                <el-tag style="margin-bottom: 5px">{{ ellipsis(val.image.split('/')[2]==undefined?val.image:val.image.split('/')[2]) }}</el-tag>
                                            </template>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="操作" width="200">
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getDaemonSetDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delDaemonSet)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="daemonset-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        :page-size="pagesize"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="daemonSetTotal">
                        </el-pagination>
                    </el-card>
                </div>
            </el-col>
        </el-row>
        <el-dialog title="YAML信息" v-model="yamlDialog" width="45%" top="5%">
            <codemirror
                :value="contentYaml"
                border
                :options="cmOptions"
                height="500"
                style="font-size:14px;"
                @change="onChange"
            ></codemirror>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="yamlDialog = false">取 消</el-button>
                    <el-button type="primary" @click="updateDaemonSet()">更 新</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script>
import common from "../common/Config";
import httpClient from '../../utils/request';
import yaml2obj from 'js-yaml';
import json2yaml from 'json2yaml';
export default {
    data() {
        return {
            cmOptions: common.cmOptions,
            contentYaml: '',
            currentPage: 1,
            pagesize: 10,
            pagesizeList: [10, 20, 30],
            searchInput: '',
            namespaceValue: 'default',
            namespaceList: [],
            namespaceListUrl: common.k8sNamespaceList,
            appLoading: false,
            daemonSetList: [],
            daemonSetTotal: 0,
            getDaemonSetsData: {
                url: common.k8sDaemonSetList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: '',
                    limit: '',
                }
            },
            daemonSetDetail: {},
            getDaemonSetDetailData: {
                url: common.k8sDaemonSetDetail,
                params: {
                    daemonset_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updateDaemonSetData: {
                url: common.k8sDaemonSetUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            delDaemonSetData: {
                url: common.k8sdaemonsetDel,
                params: {
                    daemonset_name: '',
                    namespace: '',
                }
            }
        }
    },
    methods: {
        transYaml(content) {
            return json2yaml.stringify(content)
        },
        transObj(content) {
            return yaml2obj.load(content)
        },
        onChange(val) {
            this.contentYaml = val
        },
        handleSizeChange(size) {
            this.pagesize = size;
            this.getDaemonSets()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getDaemonSets()
        },
        handleClose(done) {
            this.$confirm('确认关闭？')
            .then(() => {
                done();
            })
            .catch(() => {});
        },
        ellipsis(value) {
            return value.length>15?value.substring(0,15)+'...':value
        },
        timeTrans(timestamp) {
            let date = new Date(new Date(timestamp).getTime() + 8 * 3600 * 1000)
            date = date.toJSON();
            date = date.substring(0, 19).replace('T', ' ')
            return date 
        },
        restartTotal(e) {
            let index, sum = 0
            let containerStatuses = e.row.status.containerStatuses
            for ( index in containerStatuses) {
                sum = sum + containerStatuses[index].restartCount 
            }
            return sum
        },
        getNamespaces() {
            httpClient.get(this.namespaceListUrl)
            .then(res => {
                this.namespaceList = res.data.items
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        getDaemonSets() {
            this.appLoading = true
            this.getDaemonSetsData.params.filter_name = this.searchInput
            this.getDaemonSetsData.params.namespace = this.namespaceValue
            this.getDaemonSetsData.params.page = this.currentPage
            this.getDaemonSetsData.params.limit = this.pagesize
            httpClient.get(this.getDaemonSetsData.url, {params: this.getDaemonSetsData.params})
            .then(res => {
                this.daemonSetList = res.data.items
                this.daemonSetTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.appLoading = false
        },
        getDaemonSetDetail(e) {
            this.getDaemonSetDetailData.params.daemonset_name = e.row.metadata.name
            this.getDaemonSetDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getDaemonSetDetailData.url, {params: this.getDaemonSetDetailData.params})
            .then(res => {
                this.daemonSetDetail = res.data
                this.contentYaml = this.transYaml(this.daemonSetDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updateDaemonSet() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updateDaemonSetData.params.namespace = this.namespaceValue
            this.updateDaemonSetData.params.content = content
            httpClient.put(this.updateDaemonSetData.url, this.updateDaemonSetData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.yamlDialog = false
        },
        delDaemonSet(e) {
            this.delDaemonSetData.params.daemonset_name = e.row.metadata.name
            this.delDaemonSetData.params.namespace = this.namespaceValue
            httpClient.delete(this.delDaemonSetData.url, {data: this.delDaemonSetData.params})
            .then(res => {
                this.getDaemonSets()
                this.$message.success({
                message: res.msg
                })
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        handleConfirm(obj, operateName, fn) {
            this.confirmContent = '确认继续 ' + operateName + ' 操作吗？'
            this.$confirm(this.confirmContent,'提示',{
                confirmButtonText: '确定',
                cancelButtonText: '取消',
            })
            .then(() => {
                fn(obj)
                })
            .catch(() => {
                this.$message.info({
                    message: '已取消操作'
                })          
            })
        },
    },
    watch: {
        namespaceValue: {
            handler() {
                localStorage.setItem('namespace', this.namespaceValue)
                this.currentPage = 1
                this.getDaemonSets()
            }
        },
    },
    beforeMount() {
        if (localStorage.getItem('namespace') !== undefined && localStorage.getItem('namespace') !== null) {
            this.namespaceValue = localStorage.getItem('namespace')
        }
        this.getNamespaces()
        this.getDaemonSets()
    }
}
</script>


<style scoped>
    .daemonset-head-card,.daemonset-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .daemonset-head-search {
        width:160px;
        margin-right:10px; 
    }
    .daemonset-body-daemonsetname {
        color: #4795EE;
    }
    .daemonset-body-daemonsetname:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
</style>