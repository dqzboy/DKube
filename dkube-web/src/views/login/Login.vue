<template>
    <div class="login">
        <el-card class="login-card">
            <template #header>
                <div class="login-card-header">
                    <span>用户登录</span>
                </div>
            </template>
            <el-form :model="loginData" :rules="loginDataRules" ref="loginData">
                <el-form-item prop="username">
                    <el-input prefix-icon="UserFilled" v-model.trim="loginData.username" maxlength="32" placeholder="请输入账号" clearable></el-input>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input prefix-icon="Lock" v-model.trim="loginData.password" maxlength="16" show-password placeholder="请输入密码" clearable></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" style="width: 100%;border-radius: 2px" :loading="loginLoading" @click="handleLogin">登 录</el-button>
                </el-form-item>
            </el-form>
        </el-card>
    </div>
</template>

<script>
import common from "../common/Config";
import httpClient from '../../utils/request';
import moment from 'moment';
import jwt from 'jsonwebtoken';
export default{
    data() {
        return {
            loginLoading: false,
            loginUrl: common.loginAuth,
            loginData: {
                username: '',
                password: ''
            },
            loginDataRules: {
                username: [{
                    required: true,
                    message: '请填写用户名',
                    trigger: 'change'
                }],
                password: [{
                    required: true,
                    message: '请填写密码',
                    trigger: 'change'
                }],
            }
        }
    },
    methods: {
        handleLogin() {
            httpClient.post(this.loginUrl, this.loginData)
            .then(res => {
                localStorage.setItem('username', this.loginData.username);
                localStorage.setItem('loginDate', moment().format('YYYY-MM-DD HH:mm:ss'));
                let token = jwt.sign(this.loginData, 'devops', { expiresIn: '10h' });
                localStorage.setItem('token', token);
                this.$router.push('/');
                this.$message.success({
                    message: res.msg
                })
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        }
    }
}
</script>

<style scoped>
    .login {
        position: absolute;
        width: 100%;
        height: 100%;
        background: aquamarine;
        background-image: url(../../assets/img/login3.webp);
        background-size: 100%;
    }
    .login-card {
        position: absolute;
        left: 40%;
        top: 30%;
        width: 350px;
        border-radius: 5px;
        background: rgb(255, 255, 255);
        overflow: hidden;
    }
    .login-card-header {
        text-align: center;
    }
</style>