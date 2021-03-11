package paging

/*
   //vue部分通过定义一个paing结构体，如下
   data () {
      return {
        paging: {
          setup: {
            keyword: "", //搜索关键字
            order: "created desc",  //排序
            limit: 10,  //每页条目数
            page: 1, //当前页面号，以1开始
            meta: {} // 其他元数据
          }, //setup部分传递给el-pagination之类的组件
          result: {
            total: 0,  //一共有多少条数据
            data: [],  //实际返回数据
            error: null, //错误信息
          }
      },
   }

   //js api部分，在查询的时候传递，key和paging两个参数。如下
   export async function queryOSRepo(paging) {
     paing.result = await axios({
       url: '/osrepo/query',
       parmas: paging.setup,
     })
   }

   //go api部分，如下
   func (s *repoService) Query(ctx *gin.Context) {
        all := // 通过实际逻辑查询到结果
	var ret []PackageInfo
        //直接使用paging.WithSlice或者page.WithGORM返回结果
	ctx.JSON(http.StatusOK, paging.WithSlice(ctx, all, &ret))
    }
*/
