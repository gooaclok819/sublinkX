<script setup lang='ts'>
import { ref,onMounted  } from 'vue'
import {getSubs,AddSub,DelSub,UpdateSub} from "@/api/subcription/subs"
import {getNodes} from "@/api/subcription/node"
interface Sub {
  ID: number;
  Name: string;
  CreateDate: string;
  Config: Config;
  Nodes: Node[];
  SubLogs:SubLogs[];
}
interface Node {
  ID: number;
  Name: string;
  Link: string;
  CreateDate: string;
}
interface Config {
  clash: string;
  udp: string;
  cert: string;
}
interface SubLogs {
  date: string;
  name: string;
  count: number;
  address: string;
}

const tableData = ref<Sub[]>([])
const Clash = ref('')
const SubTitle = ref('')
const Subname = ref('')
const oldSubname = ref('')
const dialogVisible = ref(false)
const Delvisible = ref(false)
const table = ref()
const NodesList = ref<Node[]>([])
const value1 = ref<string[]>([])
const checkList = ref<string[]>([]) // 配置列表
const iplogsdialog = ref(false)
const IplogsList = ref<SubLogs[]>([])

async function getsubs() {
  const {data} = await getSubs();
    tableData.value = data
}
onMounted(() => {
    getsubs()

})
onMounted(async() => {
    const {data} = await getNodes();
    NodesList.value = data
})


const addSubs = async ()=>{
    const config = JSON.stringify({
    "clash": Clash.value.trim(),
    "udp": checkList.value.includes('udp') ? true :  false,
    "cert": checkList.value.includes('cert') ? true :  false

  })
  if (SubTitle.value === '添加订阅') {
    await AddSub({
      config: config,
      name: Subname.value.trim(),
      nodes: value1.value.join(',')
    })
    getsubs()
    ElMessage.success("添加成功");
  }else{
    await UpdateSub({
      config: config,
      name: Subname.value.trim(),
      nodes: value1.value.join(','),
      oldname: oldSubname.value
    })
    getsubs()
    ElMessage.success("更新成功");
  }

    dialogVisible.value = false;
}

const multipleSelection = ref<Sub[]>([])
const handleSelectionChange = (val: Sub[]) => {
  multipleSelection.value = val
  
}
const selectAll = () => {
  tableData.value.forEach(row => {
            table.value.toggleRowSelection(row, true)
        })
}
const handleIplogs = (row: any) => {
  iplogsdialog.value = true
  nextTick(() => {
    tableData.value.forEach((item) => {
    if (item.ID === row.ID) {
      IplogsList.value = item.SubLogs
    }
  })
  
  })
}

const toggleSelection = () => {
  table.value.clearSelection()
}

const handleAddSub = ()=>{
  SubTitle.value = '添加订阅'
  Subname.value = ''
  oldSubname.value = ''
  checkList.value = []
  Clash.value = './template/clash.yaml'
  dialogVisible.value = true
  value1.value = []
}
const handleEdit = (row:any) => {
  for (let i = 0; i < tableData.value.length; i++) {
    if (tableData.value[i].ID === row.ID) {
      function toConfig(value: string | Config): Config {
        if (typeof value === 'string') {
          return JSON.parse(value) as Config;
        } else {
          return value as Config;
        }
      }
      const config = toConfig(tableData.value[i].Config);
      SubTitle.value = '编辑订阅'
      Subname.value = tableData.value[i].Name
      oldSubname.value = Subname.value
      if (config.udp)  {
        checkList.value.push('udp')
      }
      if (config.cert)  {
        checkList.value.push('cert')
      }
      Clash.value = config.clash
      dialogVisible.value = true
      value1.value = tableData.value[i].Nodes.map((item) => item.Name)
    }
  }
}
const handleDel = (row:any) => {
  ElMessageBox.confirm(
    `你是否要删除 ${row.Name} ?`,
    '提示',
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  ).then(async () => {
      await DelSub({
        id: row.ID
      })
      getsubs()
      ElMessage({
        type: 'success',
        message: '删除成功',
      })
      
    })
}

const selectDel = () => {
  if (multipleSelection.value.length === 0) {
      return
  }
  ElMessageBox.confirm(
    `你是否要删除选中这些 ?`,
    '提示',
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  ).then( () => {
    for (let i = 0; i < multipleSelection.value.length; i++) {
       DelSub({
        id: multipleSelection.value[i].ID
      })
        tableData.value = tableData.value.filter((item) => item.ID !== multipleSelection.value[i].ID)
      }
      ElMessage({
        type: 'success',
        message: '删除成功',
      })
      
      
    })

}
// 分页显示
const currentPage = ref(1);
const pageSize = ref(10);
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  // console.log(`每页 ${val} 条`);
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val;
}
// 表格数据静态化
const currentTableData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return tableData.value.slice(start, end);
});
const isNameShow = (row: any): boolean =>  {
  return row.Name.length > 10;
}
// 复制链接
const copyInfo = (row: any) => {
  navigator.clipboard.writeText(row.Link).then(function() {
    ElMessage({
        type: 'success',
        message: '复制成功！',
      })
}, function(err) {
  ElMessage({
        type: 'warning',
        message: '复制失败！',
      })
});
}
const handleOpenUrl = (index: string,row: any) => {
  let base64data =  window.btoa(unescape(encodeURIComponent(row.Name)));
  let url = `${import.meta.env.VITE_APP_API_URL}/c/${index}/${base64data}`
  window.open(url)
}

</script>

<template>
  <div>
    <el-dialog v-model="iplogsdialog" title="访问记录" width="80%" draggable>
  <template #footer>
    <div class="dialog-footer">
      <el-table :data="IplogsList" border style="width: 100%">
        <el-table-column prop="IP" label="Ip" />
        <el-table-column prop="Count" label="当日访问次数" />
        <el-table-column prop="Addr" label="来源" />
        <el-table-column prop="Date" label="时间" />
      </el-table>
    </div>
  </template>
</el-dialog>
    <el-dialog
    v-model="dialogVisible"
    :title="SubTitle"
  >
  <el-input v-model="Subname" placeholder="请输入订阅名称" />
  
  <el-row >
  <el-tag type="primary">clash本地模版文件或url连接</el-tag>
  <el-input v-model="Clash" placeholder="clash模版文件"  />
</el-row>

  <el-row>
    <el-tag type="primary">强制开启选项</el-tag>
  <el-checkbox-group v-model="checkList"  style="margin: 5px;">
    <el-checkbox :value="'udp'">udp</el-checkbox>
    <el-checkbox :value="'cert'">跳过证书</el-checkbox>
  </el-checkbox-group>
</el-row>
  <div class="m-4">
    <p>选择已有的节点列表</p>
    <el-select
      v-model="value1"
      multiple
      placeholder="Select"
      style="width: 100%"
    >
      <el-option
        v-for="item in NodesList"
        :key="item.Name"
        :label="item.Name"
        :value="item.Name"
      />
    </el-select>
  </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="addSubs">确定</el-button>
      </div>
    </template>
  </el-dialog>
    <el-card>
    <el-button type="primary" @click="handleAddSub">添加订阅</el-button>
    <div style="margin-bottom: 10px"></div>

      <el-table ref="table" 
      :data="currentTableData" 
      style="width: 100%" 
      stripe
      @selection-change="handleSelectionChange" 
      row-key="ID" 
      :tree-props="{children: 'Nodes'}"
      >
    <el-table-column type="selection" fixed prop="ID" label="id"  />
    <el-table-column prop="Name" label="订阅名称 / 节点"  >
    <template #default="{row}">
      <el-tag :type="row.Nodes ? 'success' : 'primary'">{{row.Name}}</el-tag>
        </template>
    </el-table-column>
    <el-table-column prop="Link" label="链接" :show-overflow-tooltip="true" >
      <template #default="{row}">
        <div v-if="row.Nodes">
              <el-button @click="handleOpenUrl('v2ray',row)">V2ray</el-button>
              <el-button @click="handleOpenUrl('clash',row)">Clash</el-button>
        </div>
        </template>
      </el-table-column>
 
    <el-table-column prop="CreateDate" label="创建时间" sortable  />
    <el-table-column  label="操作" width="120">
      <template #default="scope">
        <div v-if="scope.row.Nodes">
          <el-button link type="primary" size="small" @click="handleIplogs(scope.row)">记录</el-button>
          <el-button link type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
  <el-button link type="primary" size="small" @click="handleDel(scope.row)">删除</el-button>
        </div>
        <div v-else>
          <el-button link type="primary" size="small" @click="copyInfo(scope.row)">复制</el-button>
        </div>

      </template>
    </el-table-column>
  </el-table>
  <div style="margin-top: 20px" />
    <el-button type="info" @click="selectAll()">全选</el-button>
    <el-button type="warning" @click="toggleSelection()">取消选择</el-button>
    <el-button type="danger" @click="selectDel">批量删除</el-button>
  <div style="margin-top: 20px"/>
  <el-pagination
  @size-change="handleSizeChange"
  @current-change="handleCurrentChange"
  :current-page="currentPage"
  :page-size="pageSize"
  layout="total, sizes, prev, pager, next, jumper"
  :page-sizes="[10, 20, 30, 40]"
  :total="tableData.length">
</el-pagination>

    </el-card>
  </div>
</template>

<style scoped>
.el-card{
  margin: 10px;
}
.el-input{
  margin-bottom: 10px;
}
.el-tag{
  margin: 5px;
}

</style>