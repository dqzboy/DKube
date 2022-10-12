<template>
    <div class="namespace">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="namespace-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button disabled style="border-radius:2px;" icon="Edit" type="primary">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="namespace-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getNamespaces()">搜索</el-button>
                                </div>
                            </el-col>
                            <el-col :span="2" :offset="14">
                                <div>
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getNamespaces()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="namespace-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="namespaceList"
                        v-loading="appLoading">
                            <el-table-column width="20"></el-table-column>
                            <el-table-column align=left label="Namespace名">
                                <template v-slot="scope">
                                    <a class="namespace-body-namespacename">{{ scope.row.metadata.name }}</a>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="标签" min-width='120'>
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
                            <el-table-column align=center prop="status.phase" label="状态">
                                <template v-slot="scope">
                                    <span :class="[scope.row.status.phase === 'Active' ? 'success-status' : 'error-status']">{{ scope.row.status.phase }}</span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="操作" min-width="120">
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getNamespaceDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delNamespace)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="namespace-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        :page-size="pagesize"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="namespaceTotal">
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
                    <el-button disabled type="primary" @click="updateNamespace()">更 新</el-button>
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
            namespaceTotal: 0,
            getNamespacesData: {
                url: common.k8sNamespaceList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: 1,
                    limit: 10,
                }
            },
            namespaceDetail: {},
            getNamespaceDetailData: {
                url: common.k8sNamespaceDetail,
                params: {
                    namespace_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updateNamespaceData: {
                url: common.k8sNamespaceUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            delNamespaceData: {
                url: common.k8snamespaceDel,
                params: {
                    namespace_name: '',
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
            this.getNamespaces()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getNamespaces()
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
        specTrans(str) {
            if ( str.indexOf('Ki') == -1 ) {
                return str
            }
            let num = str.slice(0,-2) / 1024 / 1024
            return num.toFixed(0)
        },
        getNamespaces() {
            this.getNamespacesData.params.filter_name = this.searchInput
            this.getNamespacesData.params.limit = this.pagesize
            this.getNamespacesData.params.page = this.currentPage
            httpClient.get(this.getNamespacesData.url, {params: this.getNamespacesData.params})
            .then(res => {
                this.namespaceList = res.data.items
                this.namespaceTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        getNamespaceDetail(e) {
            this.getNamespaceDetailData.params.namespace_name = e.row.metadata.name
            this.getNamespaceDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getNamespaceDetailData.url, {params: this.getNamespaceDetailData.params})
            .then(res => {
                this.namespaceDetail = res.data
                this.contentYaml = this.transYaml(this.namespaceDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updateNamespace() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updateNamespaceData.params.namespace = this.namespaceValue
            this.updateNamespaceData.params.content = content
            httpClient.put(this.updateNamespaceData.url, this.updateNamespaceData.params)
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
        delNamespace(e) {
            this.delNamespaceData.params.namespace_name = e.row.metadata.name
            this.delNamespaceData.params.namespace = this.namespaceValue
            httpClient.delete(this.delNamespaceData.url, {data: this.delNamespaceData.params})
            .then(res => {
                this.getNamespaces()
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
    beforeMount() {
        this.getNamespaces()
    }
}
</script>


<style scoped>
    .namespace-head-card,.namespace-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .namespace-head-search {
        width:160px;
        margin-right:10px; 
    }
    .namespace-body-namespacename {
        color: #4795EE;
    }
    .namespace-body-namespacename:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
    .success-status {
        color: rgb(27, 202, 21);
    }
    .warning-status {
        color: rgb(233, 200, 16);
    }
    .error-status {
        color: rgb(226, 23, 23);
    }
</style>