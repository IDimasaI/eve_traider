<template>
    <div
        class="flex justify-between w-full text-left bg-white shadow-xl rounded-md p-4 border-[1.5px] border-gray-300 mx-auto">
        <div class="float-left">
            <slot name="left"></slot>
        </div>
        <div>
            <slot name="center"></slot>
        </div>
        <div class="float-right">
            <slot name="right"></slot>
        </div>
    </div>
</template>