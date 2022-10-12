<template>
    <div class="secret">
        <el-row>
            <el-col :span="24">
                <div>
                    <el-card class="secret-head-card" shadow="never" :body-style="{padding:'10px'}">
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
                                    <el-button style="border-radius:2px;" icon="Refresh" plain @click="getSecrets()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="secret-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row>
                            <el-col :span="2">
                                <div>
                                    <el-button disabled style="border-radius:2px;" icon="Edit" type="primary">创建</el-button>
                                </div>
                            </el-col>
                            <el-col :span="6">
                                <div>
                                    <el-input class="secret-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius:2px;" icon="Search" type="primary" plain @click="getSecrets()">搜索</el-button>
                                </div>
                            </el-col>
                        </el-row>
                    </el-card>
                </div>
            </el-col>
            <el-col :span="24">
                <div>
                    <el-card class="secret-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table
                        style="width:100%;font-size:12px;margin-bottom:10px;"
                        :data="secretList"
                        v-loading="appLoading">
                            <el-table-column width="20"></el-table-column>
                            <el-table-column align=left label="Secret名">
                                <template v-slot="scope">
                                    <a class="secret-body-secretname">{{ scope.row.metadata.name }}</a>
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
                            <el-table-column align=center label="DATA">
                                <template v-slot="scope">
                                        <el-popover
                                            style="overflow:auto"
                                            placement="right"
                                            :width="400"
                                            trigger="click">
                                            <div style="overflow-y:auto;max-height:500px;">
                                                <span>{{ scope.row.data }}</span>
                                            </div>
                                            <template #reference>
                                                <el-icon style="font-size:18px;cursor:pointer;"><reading/></el-icon>
                                            </template>
                                        </el-popover>
                                </template>
                            </el-table-column>
                            <el-table-column align=center prop="type" min-width="100" label="类型">
                            </el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="操作" width="200">
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getSecretDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delSecret)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <el-pagination
                        class="secret-body-pagination"
                        background
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                        :current-page="currentPage"
                        :page-sizes="pagesizeList"
                        layout="total, sizes, prev, pager, next, jumper"
                        :prev-click="getSecrets"
                        :total="secretTotal">
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
                    <el-button type="primary" @click="updateSecret()">更 新</el-button>
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
            secretList: [],
            secretTotal: 0,
            getSecretsData: {
                url: common.k8sSecretList,
                params: {
                    filter_name: '',
                    namespace: '',
                    page: '',
                    limit: '',
                }
            },
            secretDetail: {},
            getSecretDetailData: {
                url: common.k8sSecretDetail,
                params: {
                    secret_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updateSecretData: {
                url: common.k8sSecretUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            delSecretData: {
                url: common.k8ssecretDel,
                params: {
                    secret_name: '',
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
            this.getSecrets()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getSecrets()
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
        getSecrets() {
            this.appLoading = true
            this.getSecretsData.params.filter_name = this.searchInput
            this.getSecretsData.params.namespace = this.namespaceValue
            this.getSecretsData.params.page = this.currentPage
            this.getSecretsData.params.limit = this.pagesize
            httpClient.get(this.getSecretsData.url, {params: this.getSecretsData.params})
            .then(res => {
                this.secretList = res.data.items
                this.secretTotal = res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.appLoading = false
        },
        getSecretDetail(e) {
            this.getSecretDetailData.params.secret_name = e.row.metadata.name
            this.getSecretDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getSecretDetailData.url, {params: this.getSecretDetailData.params})
            .then(res => {
                this.secretDetail = res.data
                this.contentYaml = this.transYaml(this.secretDetail)
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },
        updateSecret() {
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updateSecretData.params.namespace = this.namespaceValue
            this.updateSecretData.params.content = content
            httpClient.put(this.updateSecretData.url, this.updateSecretData.params)
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
        delSecret(e) {
            this.delSecretData.params.secret_name = e.row.metadata.name
            this.delSecretData.params.namespace = this.namespaceValue
            httpClient.delete(this.delSecretData.url, {data: this.delSecretData.params})
            .then(res => {
                this.getSecrets()
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
                this.getSecrets()
            }
        },
    },
    beforeMount() {
        if (localStorage.getItem('namespace') !== undefined && localStorage.getItem('namespace') !== null) {
            this.namespaceValue = localStorage.getItem('namespace')
        }
        this.getNamespaces()
        this.getSecrets()
    }
}
</script>


<style scoped>
    .secret-head-card,.secret-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .secret-head-search {
        width:160px;
        margin-right:10px; 
    }
    .secret-body-secretname {
        color: #4795EE;
    }
    .secret-body-secretname:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
</style>