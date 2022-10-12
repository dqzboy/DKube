<template>
    <div class="main-body-div">
        <el-row>
            <el-col :span="24">
                <div>
                    <img class="main-body-img" src="../../assets/img/404.png" />
                </div>
                </el-col>
            <el-col :span="24">
                <div>
                    <p class="status-code">404</p>
                    <p class="status-describe">你所访问的页面不存在······</p>
                </div>
            </el-col>
        </el-row>
    </div>
</template>

<script>
export default {
}
</script>


<style scoped>
  .main-body-img {
    margin-top: 15%;
  }
  .main-body-div {
    text-align: center;
    height: 100vh;
    width: 100vw;
  }
  .status-code {
    margin: 20px 0 20px 0;
    font-size: 95px;
    font-weight: bold;
    color: rgb(54, 95, 230);
  }
  .status-describe {
    color: rgb(145, 143, 143);
  }
</style>