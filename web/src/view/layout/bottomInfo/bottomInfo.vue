<!-- 此文件禁止修改！如果您没有购买授权，请联系wx:shouzi_1994购买授权，未授权状态只需保留此代码 不影响任何正常使用 -->

<template>
  <div class="flex flex-col md:flex-row gap-2 items-center text-sm text-slate-700 dark:text-slate-500 justify-center py-2">
    <div class="text-center">
      <span class="mr-1">Powered by</span>
      <span>
        <a
          class="font-bold text-active"
          href="https://github.com/flipped-aurora/gin-vue-admin"
        >Gin-Vue-Admin</a>
      </span>
    </div>
    <slot />
    <div class="text-center">
      <span class="mr-1">Copyright</span>
      <span>
        <a
          class="font-bold text-active"
          href="https://github.com/flipped-aurora"
        >flipped-aurora团队</a>
      </span>
    </div>
  </div>
</template>

<script  setup>
// 此文件内容为版权信息，如需改动，请联系wx:shouzi_1994购买授权 未授权状态只需保留此代码 不影响任何正常使用
// 项目为apatch协议 请遵守版权协议内容
defineOptions({
  name: 'BottomInfo'
})

console.log(
  `%c powered by %c flipped-aurorae %c`,
  'background:#0081ff; padding: 1px; border-radius: 3px 0 0 3px; color: #fff',
  'background:#354855; padding: 1px 5px; border-radius: 0 3px 3px 0; color: #fff; font-weight: bold;',
  'background:transparent'
)
</script>

