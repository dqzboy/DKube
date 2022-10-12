<template>
    <div class="statefulset">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="statefulset-head-card" shadow="never" :body-style="{padding:'10px'}">
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
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getStatefulSets()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="statefulset-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button disabled style="border-radius:2px;" icon="Edit" type="primary">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="statefulset-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getStatefulSets()">搜索</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="statefulset-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="statefulSetList"
                        v-loading="appLoading">
                            <el-table-column width="20"></el-table-column>
                            <el-table-column align=left label="StatefulSet名">
                                <template v-slot="scope">
                                    <a class="statefulset-body-statefulsetname">{{ scope.row.metadata.name }}</a>
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
                                    <span>{{ scope.row.status.currentReplicas>0?scope.row.status.currentReplicas:0  }} / {{ scope.row.spec.replicas>0?scope.row.spec.replicas:0 }} </span>
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
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getStatefulSetDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delStatefulSet)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="statefulset-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        :page-size="pagesize"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="statefulSetTotal">
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
                    <el-button type="primary" @click="updateStatefulSet()">更 新</el-button>
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
            statefulSetList: [],
            statefulSetTotal: 0,
            getStatefulSetsData: {
                url: common.k8sStatefulSetList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: '',
                    limit: '',
                }
            },
            statefulSetDetail: {},
            getStatefulSetDetailData: {
                url: common.k8sStatefulSetDetail,
                params: {
                    statefulset_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updateStatefulSetData: {
                url: common.k8sStatefulSetUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            delStatefulSetData: {
                url: common.k8sstatefulsetDel,
                params: {
                    statefulset_name: '',
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
            this.getStatefulSets()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getStatefulSets()
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
        getStatefulSets() {
            this.appLoading = true
            this.getStatefulSetsData.params.filter_name = this.searchInput
            this.getStatefulSetsData.params.namespace = this.namespaceValue
            this.getStatefulSetsData.params.page = this.currentPage
            this.getStatefulSetsData.params.limit = this.pagesize
            httpClient.get(this.getStatefulSetsData.url, {params: this.getStatefulSetsData.params})
            .then(res => {
                this.statefulSetList = res.data.items
                this.statefulSetTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.appLoading = false
        },
        getStatefulSetDetail(e) {
            this.getStatefulSetDetailData.params.statefulset_name = e.row.metadata.name
            this.getStatefulSetDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getStatefulSetDetailData.url, {params: this.getStatefulSetDetailData.params})
            .then(res => {
                this.statefulSetDetail = res.data
                this.contentYaml = this.transYaml(this.statefulSetDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updateStatefulSet() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updateStatefulSetData.params.namespace = this.namespaceValue
            this.updateStatefulSetData.params.content = content
            httpClient.put(this.updateStatefulSetData.url, this.updateStatefulSetData.params)
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
        delStatefulSet(e) {
            this.delStatefulSetData.params.statefulset_name = e.row.metadata.name
            this.delStatefulSetData.params.namespace = this.namespaceValue
            httpClient.delete(this.delStatefulSetData.url, {data: this.delStatefulSetData.params})
            .then(res => {
                this.getStatefulSets()
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
                this.getStatefulSets()
            }
        },
    },
    beforeMount() {
        if (localStorage.getItem('namespace') !== undefined && localStorage.getItem('namespace') !== null) {
            this.namespaceValue = localStorage.getItem('namespace')
        }
        this.getNamespaces()
        this.getStatefulSets()
    }
}
</script>


<style scoped>
    .statefulset-head-card,.statefulset-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .statefulset-head-search {
        width:160px;
        margin-right:10px; 
    }
    .statefulset-body-statefulsetname {
        color: #4795EE;
    }
    .statefulset-body-statefulsetname:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
</style>

